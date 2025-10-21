package mod_db

import (
	"errors"
	"net/netip"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_net"
	"rmm23/src/mod_strings"
)

func (r *RedisRepository) CheckIPHostNumber(usersSubnet netip.Prefix, userBits int) (err error) {
	var (
		entries []*Entry
	)

	switch _, entries, err = r.SearchEntryFVs(
		&mod_strings.FVs{
			{
				mod_strings.F_type,
				entryTypeUser.Number() + " " + entryTypeUser.Number(),
			},
		},
	); {
	case err != nil:
		l.Z{l.E: err}.Critical()

		return
	}

	for _, b := range entries {
		switch {
		case len(b.IPHostNumber) == 1:
			switch swErr := mod_net.Subnets.PrefixUse(usersSubnet, userBits, b.IPHostNumber[0]); {
			case errors.Is(swErr, mod_errors.EEXIST):
				l.Z{l.E: "prefix already in use", "DN": b.DN.String(), "prefix": b.IPHostNumber[0].String()}.Informational()
				b.IPHostNumber = b.IPHostNumber[:0]
				b.Status = entryStatusUpdate
				_ = r.UpdateEntry(b)
			case swErr != nil:
				l.Z{l.E: "invalid prefix", "DN": b.DN.String(), "prefix": b.IPHostNumber[0].String()}.Informational()
				b.IPHostNumber = b.IPHostNumber[:0]
				b.Status = entryStatusUpdate
				_ = r.UpdateEntry(b)
			}
		}
	}

	for _, b := range entries {
		switch {
		case len(b.IPHostNumber) > 1:
			l.Z{l.E: "prefixes are >1", "DN": b.DN.String(), "prefix": b.IPHostNumber}.Informational()

			var (
				prefixes = b.IPHostNumber
			)

			b.IPHostNumber = b.IPHostNumber[:0]
			// b.Status = entryStatusUpdate

			for _, d := range prefixes {
				switch swErr := mod_net.Subnets.PrefixUse(usersSubnet, userBits, d); {
				case errors.Is(swErr, mod_errors.EEXIST):
					l.Z{l.E: "prefix already in use", "DN": b.DN.String(), "prefix": b.IPHostNumber[0].String()}.Informational()

					continue
				case swErr != nil:
					l.Z{l.E: "invalid prefix", "DN": b.DN.String(), "prefix": b.IPHostNumber[0].String()}.Informational()

					continue
				default:
					b.IPHostNumber = []netip.Prefix{d}
					b.Status = entryStatusUpdate
					_ = r.UpdateEntry(b)
				}
			}
		}
	}

	for _, b := range entries {
		switch {
		case len(b.IPHostNumber) == 0:
			switch prefix, swErr := mod_net.Subnets.PrefixUseFree(usersSubnet, userBits); {
			case swErr != nil:
				l.Z{l.E: "no free prefix", "DN": b.DN.String()}.Warning()
			default:
				l.Z{l.M: "use new prefix", "DN": b.DN.String(), "prefix": prefix.String()}.Informational()
				b.IPHostNumber = []netip.Prefix{prefix}
				b.Status = entryStatusUpdate
				_ = r.UpdateEntry(b)
			}
		}
	}

	return
}
