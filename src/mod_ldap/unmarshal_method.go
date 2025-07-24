package mod_ldap

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

// call `Normalize` in each method instead of from `UnmarshalEntry` in hope that some time `go-ldap` will implement custom unmarshal methods.
// according to `LDAP` spec, output is not ordered.

func (r *AttrDN) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim *ldap.DN
		)

		switch interim, err = ldap.ParseDN(value); {
		case err != nil:
			continue
		}

		*r = AttrDN(interim.String())

		return
	}

	return
}

func (r *AttrDNs) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim *ldap.DN
		)

		switch interim, err = ldap.ParseDN(value); {
		case err != nil:
			continue
		}

		*r = append(*r, AttrDN(interim.String()))
	}

	return
}

func (r *AttrDestinationIndicators) UnmarshalLDAPAttr(values []string) (err error) {
	*r = mod_slices.StringsNormalize(values, mod_slices.FlagNormalize)

	return
}

func (r *AttrID) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		*r = AttrID(value)

		return // return only first value
	}

	return
}

func (r *AttrIDNumber) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim uint64
		)

		switch interim, err = strconv.ParseUint(value, 0, 0); {
		case err != nil:
			continue
		}

		*r = AttrIDNumber(interim)

		return // return only first value
	}

	return
}

func (r *AttrIPHostNumbers) UnmarshalLDAPAttr(values []string) (err error) {
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

func (r *AttrLabeledURIs) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim = strings.SplitN(value, " ", mod_slices.KVElements)
		)

		switch len(interim) {
		case 1:
			*r = append(*r, LabeledURILegacy{Key: interim[0]})
		case mod_slices.KVElements:
			*r = append(*r, LabeledURILegacy{Key: interim[0], Value: interim[1]})
		}
	}

	return
}

func (r *AttrMails) UnmarshalLDAPAttr(values []string) (err error) {
	*r = mod_slices.StringsNormalize(values, mod_slices.FlagNormalize)

	return
}

func (r *AttrObjectClasses) UnmarshalLDAPAttr(values []string) (err error) {
	*r = mod_slices.StringsNormalize(values, mod_slices.FlagNormalize)

	return
}

func (r *AttrSSHPublicKeys) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		(*r)[value] = mod_ssh.PublicKey(value)
	}

	return
}

func (r *AttrString) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		*r = AttrString(value)

		return // return only first value
	}

	return
}

func (r *AttrStrings) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		*r = append(*r, AttrString(value))
	}

	return
}

func (r *AttrTimestamp) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim time.Time
		)

		switch interim, err = ber.ParseGeneralizedTime([]byte(value)); {
		case err != nil:
			continue
		}

		*r = AttrTimestamp(interim)

		return // return only first value
	}

	return
}

func (r *AttrUserPassword) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		*r = AttrUserPassword(value)

		return // return only first value
	}

	return
}

func (r *AttrUserPKCS12s) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim *mod_crypto.Certificate
		)

		switch interim, err = mod_crypto.ParsePEM([]byte(value)); {
		case err != nil:
			continue
		}

		(*r)[AttrDN(interim.Certificates[0].Subject.String())] = *interim
	}

	return
}

func (r *AttrUUID) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim uuid.UUID
		)

		switch interim, err = uuid.Parse(value); {
		case err != nil:
			continue
		}

		*r = AttrUUID(interim)

		return // return only first value
	}

	return
}
