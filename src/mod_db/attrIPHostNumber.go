package mod_db

import (
	"errors"
	"net/netip"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_net"
)

func CheckIPHostNumber(entries []*Entry, usersSubnet netip.Prefix, userBits int) {
	for _, b := range entries {
		switch len(b.IPHostNumber) {
		case 0:
		case 1:
			switch swErr := mod_net.Subnets.PrefixUse(usersSubnet, userBits, b.IPHostNumber[0]); {
			case errors.Is(swErr, mod_errors.EEXIST):
				l.Z{l.E: "prefix already in use", "DN": b.DN.String(), "prefix": b.IPHostNumber[0].String()}.Informational()
				b.IPHostNumber = b.IPHostNumber[:0]
				b.Ver++
			case swErr != nil:
				l.Z{l.E: "invalid prefix", "DN": b.DN.String(), "prefix": b.IPHostNumber[0].String()}.Informational()
				b.IPHostNumber = b.IPHostNumber[:0]
				b.Ver++
			}
		default:
		}
	}

	for _, b := range entries {
		switch len(b.IPHostNumber) {
		case 0:
		case 1:
		default:
			l.Z{l.E: "prefixes are >1", "DN": b.DN.String(), "prefix": b.IPHostNumber}.Informational()
		}
	}

	for _, b := range entries {
		switch len(b.IPHostNumber) {
		case 0:
			switch prefix, swErr := mod_net.Subnets.PrefixUseFree(usersSubnet, userBits); {
			case swErr != nil:
				l.Z{l.E: "no free prefix", "DN": b.DN.String()}.Warning()
			default:
				l.Z{l.M: "use new prefix", "DN": b.DN.String(), "prefix": prefix.String()}.Informational()
				b.IPHostNumber = append(b.IPHostNumber, prefix)
				b.Ver++
			}
		case 1:
		default:
		}
	}
}
