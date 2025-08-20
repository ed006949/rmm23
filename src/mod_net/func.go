package mod_net

import (
	"errors"
	"net"
	"net/netip"

	"rmm23/src/mod_errors"
)

func LookupMX(names []string) (outbound []string, errs mod_errors.Errs) {
	for _, name := range names {
		var (
			interim, err = net.LookupMX(name)
			errDetail    *net.DNSError
			_            = errors.As(err, &errDetail)
		)
		switch {
		case errDetail != nil && errDetail.IsNotFound:
			continue
		case err != nil:
			errs = append(errs, err)

			continue
		}

		for _, b := range interim {
			outbound = append(outbound, b.Host)
		}
	}

	return
}

// func UrlParse(inbound string) (outbound *url.URL, err error) {
// 	switch outbound, err = url.Parse(inbound); {
// 	case err != nil:
// 		return nil, err
// 	case len(outbound.String()) == 0:
// 		return nil, mod_errors.ENODATA
// 	}
//
// 	return
// }

func ParseNetIP(inbound net.IP) (outbound *netip.Addr, err error) {
	switch ip4 := inbound.To4(); {
	case ip4 != nil:
		var (
			arr [4]byte
			ip  netip.Addr
		)
		copy(arr[:], ip4)
		ip = netip.AddrFrom4(arr)

		return &ip, nil
	}

	switch ip16 := inbound.To16(); {
	case ip16 != nil:
		var (
			arr [16]byte
			ip  netip.Addr
		)
		copy(arr[:], ip16)
		ip = netip.AddrFrom16(arr)

		return &ip, nil
	}

	return nil, mod_errors.ENODATA
}

func ParseNetIPs(inbound []net.IP) (outbound []*netip.Addr, err error) {
	outbound = make([]*netip.Addr, len(inbound), len(inbound))
	for a, b := range inbound {
		outbound[a] = mod_errors.StripErr1(ParseNetIP(b))
	}

	return
}

// func CleanPath(inbound *url.URL) (outbound string) { return strings.TrimPrefix(inbound.Path, "/") }

// func CleanUser(inbound *url.URL) (username string, password string) {
// 	return CleanUsername(inbound), CleanPassword(inbound)
// }

// func CleanUsername(inbound *url.URL) (outbound string) { return inbound.User.Username() }
// func CleanPassword(inbound *url.URL) (outbound string) {
// 	return mod_bools.StripIfBool1(inbound.User.Password())
// }

// func RedisNetwork(inbound *url.URL) (outbound string, err error) {
// 	switch outbound = inbound.Scheme; outbound {
// 	case "redis", "redis-sentinel":
// 		return "tcp", nil
// 	case "file":
// 		return "unix", nil
// 	default:
// 		return outbound, mod_errors.EUnknownScheme
// 	}
// }
