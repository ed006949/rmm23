package mod_db

import (
	"net/netip"
	"strconv"
	"strings"
	"time"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/mod_crypto"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_ssh"
)

// call `Normalize` from each method instead of from `UnmarshalEntry` in hope that sometime `go-ldap` will implement custom marshal/unmarshal mechanics.
// according to `LDAP` spec, output is not ordered.

func (r *attrDN) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim *ldap.DN
		)

		switch interim, err = ldap.ParseDN(value); {
		case err != nil:
			continue
		}

		*r = attrDN(interim.String())

		return
	}

	return
}

func (r *attrDNs) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim *ldap.DN
		)

		switch interim, err = ldap.ParseDN(value); {
		case err != nil:
			continue
		}

		*r = append(*r, attrDN(interim.String()))
	}

	return
}

func (r *attrDestinationIndicators) UnmarshalLDAPAttr(values []string) (err error) {
	*r = mod_slices.StringsNormalize(values, mod_slices.FlagNormalize)

	return
}

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

func (r *attrIPHostNumbers) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim netip.Prefix
		)

		switch interim, err = netip.ParsePrefix(value); {
		case err != nil:
			continue
		}

		*r = append(*r, interim)
	}

	return nil
}

// func (r *attrLabeledURIs) UnmarshalLDAPAttr(values []string) (err error) {
//	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
//		var (
//			interim = strings.SplitN(value, " ", mod_slices.KVElements)
//		)
//
//		switch len(interim) {
//		case 1:
//			*r = append(*r, labeledURILegacy{Key: interim[0]})
//		case mod_slices.KVElements:
//			*r = append(*r, labeledURILegacy{Key: interim[0], Value: interim[1]})
//		}
//	}
//
//	return
// }

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

func (r *attrObjectClasses) UnmarshalLDAPAttr(values []string) (err error) {
	*r = mod_slices.StringsNormalize(values, mod_slices.FlagNormalize)

	return
}

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
		*r = append(*r, value)
	}

	return
}

func (r *attrTime) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim time.Time
		)

		switch interim, err = ber.ParseGeneralizedTime([]byte(value)); {
		case err != nil:
			continue
		}

		*r = attrTime(interim)

		return // return only first value
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

func (r *attrUserPKCS12s) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			forErr  error
			interim *mod_crypto.Certificate
		)

		switch interim, forErr = mod_crypto.ParsePEM([]byte(value)); {
		case forErr != nil:
			continue
		}

		(*r)[attrDN(interim.Certificates[0].Subject.String())] = *interim
	}

	return
}

func (r *attrUUID) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim uuid.UUID
		)

		switch interim, err = uuid.Parse(value); {
		case err != nil:
			continue
		}

		*r = attrUUID(interim)

		return // return only first value
	}

	return
}
