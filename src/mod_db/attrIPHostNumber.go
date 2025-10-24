package mod_db

import (
	"errors"
	"net/netip"

	"github.com/rs/zerolog/log"

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
		log.Error().Err(err).Send()

		return
	}

	for _, b := range entries {
		switch {
		case len(b.IPHostNumber) == 1:
			switch swErr := mod_net.Subnets.PrefixUse(usersSubnet, userBits, b.IPHostNumber[0]); {
			case errors.Is(swErr, mod_errors.EEXIST):
				log.Info().Str("DN", b.DN.String()).Str("prefix", b.IPHostNumber[0].String()).Msg("prefix already in use")
				b.IPHostNumber = b.IPHostNumber[:0]
				b.Status = entryStatusUpdate
				_ = r.UpdateEntry(b)
			case swErr != nil:
				log.Info().Str("DN", b.DN.String()).Str("prefix", b.IPHostNumber[0].String()).Msg("invalid prefix")
				b.IPHostNumber = b.IPHostNumber[:0]
				b.Status = entryStatusUpdate
				_ = r.UpdateEntry(b)
			}
		}
	}

	for _, b := range entries {
		switch {
		case len(b.IPHostNumber) > 1:
			log.Info().Str("DN", b.DN.String()).Any("prefix", b.IPHostNumber).Msg("prefixes are >1")

			var (
				prefixes = b.IPHostNumber
			)

			b.IPHostNumber = b.IPHostNumber[:0]
			// b.Status = entryStatusUpdate

			for _, d := range prefixes {
				switch swErr := mod_net.Subnets.PrefixUse(usersSubnet, userBits, d); {
				case errors.Is(swErr, mod_errors.EEXIST):
					log.Info().Str("DN", b.DN.String()).Str("prefix", b.IPHostNumber[0].String()).Msg("prefix already in use")

					continue
				case swErr != nil:
					log.Info().Str("DN", b.DN.String()).Str("prefix", b.IPHostNumber[0].String()).Msg("invalid prefix")

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
				log.Warn().Str("DN", b.DN.String()).Msg("no free prefix")
			default:
				log.Info().Str("DN", b.DN.String()).Str("prefix", prefix.String()).Msg("use new prefix")
				b.IPHostNumber = []netip.Prefix{prefix}
				b.Status = entryStatusUpdate
				_ = r.UpdateEntry(b)
			}
		}
	}

	return
}
