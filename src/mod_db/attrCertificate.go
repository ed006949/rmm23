package mod_db

import (
	"github.com/google/uuid"

	"rmm23/src/mod_bools"
	"rmm23/src/mod_crypto"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_net"
)

func (r *Cert) ParseDERs(key, crt, ca, crl, csr []byte) (err error) {
	var (
		cert = new(mod_crypto.Certificate)
	)

	switch err = cert.ParseDERs(key, crt, ca, crl, csr); {
	case err != nil:
		return
	}

	*r = Cert{
		// Key:            "",
		// Ver:            0,
		Ext:            cert.Certificate.NotAfter,
		UUID:           uuid.NewSHA1(uuid.NameSpaceOID, cert.Certificate.Raw),
		SerialNumber:   cert.Certificate.SerialNumber,
		Issuer:         mod_errors.StripErr1(parseDN(cert.Certificate.Issuer.String())),
		Subject:        mod_errors.StripErr1(parseDN(cert.Certificate.Subject.String())),
		NotBefore:      cert.Certificate.NotBefore,
		NotAfter:       cert.Certificate.NotAfter,
		DNSNames:       cert.Certificate.DNSNames,
		EmailAddresses: cert.Certificate.EmailAddresses,
		IPAddresses:    mod_errors.StripErr1(mod_net.ParseNetIPs(cert.Certificate.IPAddresses)),
		URIs:           cert.Certificate.URIs,
		IsCA:           mod_bools.AttrBool(cert.Certificate.IsCA),
		Certificate:    cert,
	}

	r.Key = r.UUID.String()

	return
}
