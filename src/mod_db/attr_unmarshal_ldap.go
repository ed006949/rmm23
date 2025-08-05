package mod_db

import (
	"net/netip"
	"strconv"
	"strings"

	"rmm23/src/mod_slices"
	"rmm23/src/mod_ssh"
)

// call `Normalize` from each method instead of from `UnmarshalEntry` in hope that sometime `go-ldap` will implement custom marshal/unmarshal mechanics.
// according to `LDAP` spec, output is not ordered.

func (r *attrDNs) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim = new(attrDN)
		)
		switch err = interim.Parse(value); {
		case err != nil:
			continue
		}

		*r = append(*r, interim)
	}

	return
}

// func (r *attrDestinationIndicators) UnmarshalLDAPAttr(values []string) (err error) {
// 	*r = mod_slices.StringsNormalize(values, mod_slices.FlagNormalize)
//
// 	return
// }

func (r *attrID) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		*r = attrID(value)

		return // return only first value
	}

	return
}

func (r *attrIDNumber) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim uint64
		)
		switch interim, err = strconv.ParseUint(value, 0, 0); {
		case err != nil:
			continue
		}

		*r = attrIDNumber(interim)

		return // return only first value
	}

	return
}

func (r *attrIPAddress) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim netip.Addr
		)
		switch interim, err = netip.ParseAddr(value); {
		case err != nil:
			continue
		}

		r.Addr = interim

		return // return only first value
	}

	return nil
}

func (r *attrIPAddresses) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim netip.Addr
		)
		switch interim, err = netip.ParseAddr(value); {
		case err != nil:
			continue
		}

		*r = append(*r, &attrIPAddress{interim})
	}

	return nil
}

func (r *attrIPPrefix) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim netip.Prefix
		)
		switch interim, err = netip.ParsePrefix(value); {
		case err != nil:
			continue
		}

		r.Prefix = interim

		return // return only first value
	}

	return nil
}

func (r *attrIPPrefixes) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim netip.Prefix
		)
		switch interim, err = netip.ParsePrefix(value); {
		case err != nil:
			continue
		}

		*r = append(*r, &attrIPPrefix{interim})
	}

	return nil
}

func (r *attrLabeledURIs) UnmarshalLDAPAttr(values []string) (err error) {
	switch {
	case *r == nil:
		*r = attrLabeledURIs{}
	}

	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim = strings.SplitN(value, " ", mod_slices.KVElements)
		)
		switch len(interim) {
		case 1:
			(*r)[interim[0]] = ""
		case mod_slices.KVElements:
			(*r)[interim[0]] = interim[1]
		}
	}

	return
}

func (r *attrMails) UnmarshalLDAPAttr(values []string) (err error) {
	*r = mod_slices.StringsNormalize(values, mod_slices.FlagNormalize)

	return
}

// func (r *attrObjectClasses) UnmarshalLDAPAttr(values []string) (err error) {
// 	*r = mod_slices.StringsNormalize(values, mod_slices.FlagNormalize)
//
// 	return
// }

func (r *attrSSHPublicKeys) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		(*r)[value] = mod_ssh.PublicKey(value)
	}

	return
}

func (r *attrString) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		*r = attrString(value)

		return // return only first value
	}

	return
}

func (r *attrStrings) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		*r = append(*r, attrString(value))
	}

	return
}

func (r *attrUserPassword) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		*r = attrUserPassword(value)

		return // return only first value
	}

	return
}
