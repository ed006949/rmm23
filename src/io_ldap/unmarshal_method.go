package io_ldap

import (
	"encoding/xml"
	"net/netip"
	"strconv"
	"time"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/io_crypto"
	"rmm23/src/io_ssh"
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
func (r *AttrUserPassword) UnmarshalLDAPAttr(values []string) (err error) {
	*r = AttrUserPassword(values[0])
	return
}
func (r *AttrID) UnmarshalLDAPAttr(values []string) (err error) {
	*r = AttrID(values[0])
	return
}
func (r *AttrString) UnmarshalLDAPAttr(values []string) (err error) {
	*r = AttrString(values[0])
	return
}

func (r AttrDNs) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		var (
			interim AttrDN
		)
		switch err = interim.UnmarshalLDAPAttr([]string{value}); {
		case err != nil:
			return
		}
		r[interim] = struct{}{}
	}
	return
}
func (r AttrDestinationIndicators) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		r[value] = struct{}{}
	}
	return
}
func (r AttrIPHostNumbers) UnmarshalLDAPAttr(values []string) (err error) {
	switch len(values) {
	case 0:
		return
	case 1:
	default:
		r.modified = true
	}
	r.data, r.invalid = netip.ParsePrefix(values[0])
	r.modified = r.modified == true || r.invalid != nil
	return
}
func (r *AttrLabeledURIs) UnmarshalLDAPAttr(values []string) (err error) {
	switch len(values) {
	case 0:
		return
	case 1:
	default:
		r.modified = true
	}
	r.invalid = xml.Unmarshal([]byte(values[0]), &r.data)
	r.modified = r.modified == true || r.invalid != nil
	return
}
func (r AttrMails) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		r[value] = struct{}{}
	}
	return
}
func (r AttrSSHPublicKeys) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		r[value] = io_ssh.PublicKey(value)
	}
	return
}
func (r AttrStrings) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		r[AttrString(value)] = struct{}{}
	}
	return
}
func (r AttrUserPKCS12s) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		var (
			interim *io_crypto.Certificate
		)
		switch interim, err = io_crypto.ParsePEM([]byte(value)); {
		case err != nil:
			continue
		}
		r[AttrDN(interim.Certificates[0].Subject.String())] = interim
	}
	return
}

// func (r *LabeledURI) UnmarshalLDAPAttr(values []string) (err error) {
// 	for _, value := range values {
// 		var (
// 			interim = strings.SplitN(value, " ", 2)
// 		)
// 		switch len(interim) {
// 		case 0:
// 		case 1:
// 		case 2:
// 		}
// 	}
// 	return
// }
