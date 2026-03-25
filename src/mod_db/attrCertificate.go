package mod_db

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"rmm23/src/mod_crypto"
	"rmm23/src/mod_dn"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_net"
	"rmm23/src/mod_time"
)

func (r *Cert) parseCertificate(inbound *mod_crypto.Certificate) (err error) {
	switch {
	case inbound == nil:
		return mod_errors.ENODATA
	}

	r.Certificate = inbound
	r.normalize()

	return
}

func (r *Cert) parseRaw(inbound ...[]byte) (err error) {
	var (
		certificate = new(mod_crypto.Certificate)
	)

	switch err = certificate.ParseRaw(inbound...); {
	case err != nil:
		log.Error().Err(err).Send()

		return
	}

	switch err = r.parseCertificate(certificate); {
	case err != nil:
		log.Error().Err(err).Send()

		return
	}

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
	r.Issuer = mod_errors.StripErr1(mod_dn.UnmarshalText([]byte(r.Certificate.Certificate.Issuer.String())))
	r.Subject = mod_errors.StripErr1(mod_dn.UnmarshalText([]byte(r.Certificate.Certificate.Subject.String())))
	r.NotBefore = mod_time.Time{r.Certificate.Certificate.NotBefore}
	r.NotAfter = mod_time.Time{r.Certificate.Certificate.NotAfter}
	r.DNSNames = r.Certificate.Certificate.DNSNames
	r.EmailAddresses = r.Certificate.Certificate.EmailAddresses
	r.IPAddresses = mod_errors.StripErr1(mod_net.ParseNetIPs(r.Certificate.Certificate.IPAddresses))
	r.URIs = r.Certificate.Certificate.URIs
	r.IsCA = r.Certificate.Certificate.IsCA
	// r.NotBeforeUnix.Set(r.NotBefore)
	// r.NotAfterUnix.Set(r.NotAfter)
}
