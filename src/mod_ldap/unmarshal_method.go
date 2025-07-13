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
	"rmm23/src/mod_ssh"
)

func (r *AttrDN) UnmarshalLDAPAttr(values []string) (err error) {
	var (
		value *ldap.DN
	)
	switch value, err = ldap.ParseDN(values[0]); {
	case err != nil:
		return
	}
	*r = AttrDN(value.String())
	return
}

func (r *AttrDNs) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		var (
			interim AttrDN
		)
		switch err = interim.UnmarshalLDAPAttr([]string{value}); {
		case err != nil:
			return
		}
		(*r)[interim] = struct{}{}
	}
	return
}

func (r *AttrDestinationIndicators) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		(*r)[value] = struct{}{}
	}
	return
}

func (r *AttrID) UnmarshalLDAPAttr(values []string) (err error) {
	*r = AttrID(values[0])
	return
}

func (r *AttrIDNumber) UnmarshalLDAPAttr(values []string) (err error) {
	var (
		value uint64
	)
	switch value, err = strconv.ParseUint(values[0], 0, 0); {
	case err != nil:
		return
	}
	*r = AttrIDNumber(value)
	return
}

func (r *AttrIPHostNumbers) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		var (
			interim netip.Prefix
		)
		switch interim, err = netip.ParsePrefix(value); {
		case err == nil:
			(*r)[interim] = struct{}{}
			return
		}
	}
	return nil
}

func (r *AttrLabeledURIs) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		var (
			interim = strings.SplitN(value, " ", 2)
		)
		switch len(interim) {
		case 0:
			continue
		case 1:
			*r = append(*r, LabeledURILegacy{Key: interim[0]})
		case 2:
			*r = append(*r, LabeledURILegacy{Key: interim[0], Value: interim[1]})
		}
	}

	return nil
}

func (r *AttrMails) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		(*r)[value] = struct{}{}
	}
	return
}

func (r *AttrObjectClasses) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		(*r)[value] = struct{}{}
	}
	return
}

func (r *AttrSSHPublicKeys) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		(*r)[value] = mod_ssh.PublicKey(value)
	}
	return
}

func (r *AttrString) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		switch interim := strings.TrimSpace(value); {
		case len(interim) > 0:
			*r = AttrString(interim)
		}
		return // return only first value
	}
	return
}

func (r *AttrStrings) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		switch interim := strings.TrimSpace(value); {
		case len(interim) > 0:
			(*r)[AttrString(interim)] = struct{}{}
		}
	}
	return
}

func (r *AttrTimestamp) UnmarshalLDAPAttr(values []string) (err error) {
	var (
		value time.Time
	)
	switch value, err = ber.ParseGeneralizedTime([]byte(values[0])); {
	case err != nil:
		return
	}
	*r = AttrTimestamp(value)
	return
}

func (r *AttrUserPassword) UnmarshalLDAPAttr(values []string) (err error) {
	*r = AttrUserPassword(values[0])
	return
}

func (r *AttrUserPKCS12s) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
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
	var (
		value uuid.UUID
	)
	switch value, err = uuid.Parse(values[0]); {
	case err != nil:
		return
	}
	*r = AttrUUID(value)
	return
}

/*// UnmarshalLDAPAttr for AttrLabeledURIs
// there must be only one valid XML data or nothing
// in other cases r.modified needs to be set, to update data in LDAP later
func (r *AttrLabeledURIs) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		var (
			data LabeledURI
		)
		switch err = xml.Unmarshal([]byte(value), &data); {
		case err == nil:
			switch {
			case r.data == nil:
				r.data = &data
			default:
				// another XML data in values - append (!good)
				r.modified = true
				r.data.OpenVPN = append(r.data.OpenVPN, data.OpenVPN...)
				r.data.CiscoVPN = append(r.data.CiscoVPN, data.CiscoVPN...)
				r.data.InterimHost = append(r.data.InterimHost, data.InterimHost...)
				r.data.Legacy = append(r.data.Legacy, data.Legacy...)
			}
			continue
		}

		// fallback to legacy key-value space-separated schema - append if any (!good)
		r.modified = true
		var (
			legacy = strings.SplitN(value, " ", 2)
		)
		switch len(legacy) {
		case 0:
			continue
		case 1:
			switch {
			case r.data == nil:
				r.data = &LabeledURI{}
			}
			r.data.Legacy = append(r.data.Legacy, LabeledURILegacy{Key: legacy[0], Value: ""})
		case 2:
			switch {
			case r.data == nil:
				r.data = &LabeledURI{}
			}
			r.data.Legacy = append(r.data.Legacy, LabeledURILegacy{Key: legacy[0], Value: legacy[1]})
		}
	}

	return nil
}
*/
