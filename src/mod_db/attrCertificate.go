package mod_db

import (
	"github.com/google/uuid"

	"rmm23/src/mod_bools"
	"rmm23/src/mod_crypto"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_net"
)

func (r *Cert) UnmarshalText(inbound []byte) (err error) {
	var (
		cert = new(Cert)
	)

	switch err = cert.parseRaw(inbound); {
	case err != nil:
		return
	}

	*r = *cert

	return
}

func (r *Cert) parseRaw(inbound ...[]byte) (err error) {
	var (
		certificate = new(mod_crypto.Certificate)
	)

	switch err = certificate.ParseRaw(inbound...); {
	case err != nil:
		return
	}

	r.normalize()

	return
}
func (r *Cert) normalize() {
	var (
		certUUID = uuid.NewSHA1(uuid.NameSpaceOID, r.Certificate.Certificate.Raw)
	)

	r.Key = certUUID.String()
	r.Ext = r.Certificate.Certificate.NotAfter
	r.UUID = certUUID
	r.SerialNumber = r.Certificate.Certificate.SerialNumber
	r.Issuer = mod_errors.StripErr1(parseDN(r.Certificate.Certificate.Issuer.String()))
	r.Subject = mod_errors.StripErr1(parseDN(r.Certificate.Certificate.Subject.String()))
	r.NotBefore = attrTime{r.Certificate.Certificate.NotBefore}
	r.NotAfter = attrTime{r.Certificate.Certificate.NotAfter}
	r.DNSNames = r.Certificate.Certificate.DNSNames
	r.EmailAddresses = r.Certificate.Certificate.EmailAddresses
	r.IPAddresses = mod_errors.StripErr1(mod_net.ParseNetIPs(r.Certificate.Certificate.IPAddresses))
	r.URIs = r.Certificate.Certificate.URIs
	r.IsCA = mod_bools.AttrBool(r.Certificate.Certificate.IsCA)
}
