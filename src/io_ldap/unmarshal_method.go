package io_ldap

import (
	"encoding/xml"
	"net/netip"
	"strconv"
	"strings"
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
	for _, value := range values {
		switch r.data, r.invalid = netip.ParsePrefix(value); {
		case r.invalid != nil:
			r.modified = true
			continue
		}
		break
	}
	r.modified = r.modified || r.invalid != nil || len(values) > 1
	return
}
func (r *AttrLabeledURIs) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range values {
		switch r.invalid = xml.Unmarshal([]byte(value), &r.data); {
		case r.invalid != nil:
			r.modified = true
			switch interim := strings.SplitN(value, " ", 2); len(interim) {
			case 0:
			case 1:
				r.data.Legacy = append(r.data.Legacy, LabeledURILegacy{Key: interim[0], Value: ""})
			case 2:
				r.data.Legacy = append(r.data.Legacy, LabeledURILegacy{Key: interim[0], Value: interim[1]})
			}
			continue
		}
		break
	}
	r.modified = r.modified || r.invalid != nil || len(values) > 1
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
