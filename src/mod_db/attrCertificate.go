package mod_db

import (
	"github.com/google/uuid"

	"rmm23/src/mod_bools"
	"rmm23/src/mod_crypto"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_net"
)

// func (r *Cert) MarshalText() (outbound []byte, err error) {
// 	var (
// 		cert = new(Cert)
// 	)
//
// 	switch err = cert.parseRaw(outbound); {
// 	case err != nil:
// 		return
// 	}
//
// 	*r = *cert
//
// 	return
// }

// func (r *Cert) UnmarshalText(inbound []byte) (err error) {
// 	var (
// 		cert = new(Cert)
// 	)
//
// 	switch err = cert.parseRaw(inbound); {
// 	case err != nil:
// 		return
// 	}
//
// 	*r = *cert
//
// 	return
// }

func (r *Cert) parseRaw(inbound ...[]byte) (err error) {
	var (
		certificate = new(mod_crypto.Certificate)
	)

	switch err = certificate.ParseRaw(inbound...); {
	case err != nil:
		return
	}

	var (
		certUUID = uuid.NewSHA1(uuid.NameSpaceOID, certificate.Certificate.Raw)
	)

	*r = Cert{
		Key:            certUUID.String(),
		Ext:            certificate.Certificate.NotAfter,
		UUID:           certUUID,
		SerialNumber:   certificate.Certificate.SerialNumber,
		Issuer:         mod_errors.StripErr1(parseDN(certificate.Certificate.Issuer.String())),
		Subject:        mod_errors.StripErr1(parseDN(certificate.Certificate.Subject.String())),
		NotBefore:      certificate.Certificate.NotBefore,
		NotAfter:       certificate.Certificate.NotAfter,
		DNSNames:       certificate.Certificate.DNSNames,
		EmailAddresses: certificate.Certificate.EmailAddresses,
		IPAddresses:    mod_errors.StripErr1(mod_net.ParseNetIPs(certificate.Certificate.IPAddresses)),
		URIs:           certificate.Certificate.URIs,
		IsCA:           mod_bools.AttrBool(certificate.Certificate.IsCA),
		Certificate:    certificate,
	}

	return
}
