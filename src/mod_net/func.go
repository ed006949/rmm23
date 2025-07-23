package mod_net

import (
	"errors"
	"net"
	"net/url"

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

func UrlParse(inbound string) (outbound *url.URL, err error) {
	switch outbound, err = url.Parse(inbound); {
	case err != nil:
		return nil, err
	case len(outbound.String()) == 0:
		return nil, mod_errors.ENODATA
	}
	return
}
