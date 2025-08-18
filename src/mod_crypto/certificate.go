package mod_crypto

import (
	"bytes"
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"software.sslmate.com/src/go-pkcs12"

	"rmm23/src/mod_errors"
)

type Certificate struct {
	P12 []byte `json:"p12,omitempty"`
	DER []byte `json:"-"` // unwilling to perform
	PEM []byte `json:"pem,omitempty" ldap:"userPKCS12"`
	CSR []byte `json:"csr,omitempty"`
	CRL []byte `json:"crl,omitempty"`

	PrivateKeyDER         []byte   `json:"-"`
	CertificateDER        []byte   `json:"-"`
	CertificateCAChainDER [][]byte `json:"-"`
	CertificateRequestDER []byte   `json:"-"`
	RevocationListDER     []byte   `json:"-"`

	PrivateKeyPEM         []byte   `json:"-"`
	CertificatePEM        []byte   `json:"-"`
	CertificateCAChainPEM [][]byte `json:"-"`
	CertificateRequestPEM []byte   `json:"-"`
	RevocationListPEM     []byte   `json:"-"`

	PrivateKey         crypto.PrivateKey        `json:"-"`
	Certificate        *x509.Certificate        `json:"-"`
	CertificateCAChain []*x509.Certificate      `json:"-"`
	CertificateRequest *x509.CertificateRequest `json:"-"`
	RevocationList     *x509.RevocationList     `json:"-"`
}

func (r *Certificate) MarshalText() (outbound []byte, err error) {
	switch {
	case r == nil:
		return nil, mod_errors.ENODATA
	case r.PEM == nil:
		switch err = r.encodePEM(); {
		case err != nil:
			return
		}
	}

	return r.PEM, nil
}

func (r *Certificate) UnmarshalText(inbound []byte) (err error) {
	var (
		interim = new(Certificate)
	)
	switch err = interim.ParseRaw(inbound); {
	case err != nil:
		return
	}

	*r = *interim

	return
}

// ParseRaw parses any certificate data.
//
// * for PEM and DER: place the certificate before the chain!
func (r *Certificate) ParseRaw(inbound ...[]byte) (err error) {
	var (
		interim = new(Certificate)
	)
	switch err = interim.parsePEM(inbound...); {
	case err == nil:
		*r = *interim

		return
	}

	switch err = interim.parseDER(inbound...); {
	case err == nil:
		*r = *interim

		return
	}

	switch {
	case len(inbound) > 0:
		switch err = interim.parseP12(inbound[0]); {
		case err == nil:
			*r = *interim

			return
		}
	}

	return mod_errors.EParse
}

// checkPrivateKey
// We don't need to parse the public key for TLS, but we so do anyway to check that it looks sane and matches the private key.
func (r *Certificate) checkPrivateKey() (err error) {
	switch {
	case r == nil || r.Certificate == nil:
		return mod_errors.ENODATA
	}

	switch pub := r.Certificate.PublicKey.(type) {
	case *rsa.PublicKey:
		return r.checkPrivateKeyRSA(pub)
	case *ecdsa.PublicKey:
		return r.checkPrivateKeyECDSA(pub)
	case ed25519.PublicKey:
		return r.checkPrivateKeyED25519(pub)
	case *ecdh.PublicKey:
		return r.checkPrivateKeyECDH(pub)
	default:
		return mod_errors.EUnknownPubKeyAlgo
	}
}
func (r *Certificate) checkPrivateKeyRSA(pub *rsa.PublicKey) (err error) {
	switch priv, ok := r.PrivateKey.(*rsa.PrivateKey); {
	case !ok:
		return mod_errors.ETypeMismatchPrivKeyPubKey
	case pub.N.Cmp(priv.N) != 0:
		return mod_errors.EMismatchPrivKeyPubKey
	}

	return
}
func (r *Certificate) checkPrivateKeyECDSA(pub *ecdsa.PublicKey) (err error) {
	switch priv, ok := r.PrivateKey.(*ecdsa.PrivateKey); {
	case !ok:
		return mod_errors.ETypeMismatchPrivKeyPubKey
	case pub.X.Cmp(priv.X) != 0 || pub.Y.Cmp(priv.Y) != 0:
		return mod_errors.EMismatchPrivKeyPubKey
	}

	return
}
func (r *Certificate) checkPrivateKeyED25519(pub ed25519.PublicKey) (err error) {
	switch priv, ok := r.PrivateKey.(ed25519.PrivateKey); {
	case !ok:
		return mod_errors.ETypeMismatchPrivKeyPubKey
	case !bytes.Equal(priv.Public().(ed25519.PublicKey), pub):
		return mod_errors.EMismatchPrivKeyPubKey
	}

	return
}
func (r *Certificate) checkPrivateKeyECDH(pub *ecdh.PublicKey) (err error) {
	switch priv, ok := r.PrivateKey.(*ecdh.PrivateKey); {
	case !ok:
		return mod_errors.ETypeMismatchPrivKeyPubKey
	case priv.PublicKey().Curve() != pub.Curve():
		return mod_errors.EMismatchPrivKeyPubKey
	case !priv.PublicKey().Equal(pub):
		return mod_errors.EMismatchPrivKeyPubKey
	}

	return
}

func (r *Certificate) parsePrivateKeyDER(der []byte) (err error) {
	var (
		key crypto.PrivateKey
	)
	switch key, err = x509.ParsePKCS1PrivateKey(der); {
	case err == nil:
		r.PrivateKey = key

		return
	}

	switch key, err = x509.ParsePKCS8PrivateKey(der); {
	case err == nil:
		switch value := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey, *ecdh.PrivateKey:
			r.PrivateKey = value

			return
		default:
			return mod_errors.EUnknownPrivKeyType
		}
	}

	switch key, err = x509.ParseECPrivateKey(der); {
	case err == nil:
		r.PrivateKey = key

		return
	}

	return mod_errors.EX509ParsePrivKey
}

func (r *Certificate) parseDER(inbound ...[]byte) (err error) {
	var (
		interim = new(Certificate)
	)

	for _, b := range inbound {
		switch err = interim.parsePrivateKeyDER(b); {
		case err == nil:
			continue
		}

		// ensure that certificate goes before CA chain
		switch interim.Certificate {
		case nil:
			switch interim.Certificate, err = x509.ParseCertificate(b); {
			case err == nil:
				continue
			}
		default:
			switch interim.CertificateCAChain, err = x509.ParseCertificates(b); {
			case err == nil:
				continue
			}
		}

		switch interim.CertificateRequest, err = x509.ParseCertificateRequest(b); {
		case err == nil:
			continue
		}

		switch interim.RevocationList, err = x509.ParseRevocationList(b); {
		case err == nil:
			continue
		}
	}

	switch err = interim.checkPrivateKey(); {
	case err != nil:
		return
	}

	*r = *interim

	return
}

func (r *Certificate) parsePEM(inbound ...[]byte) (err error) {
	var (
		interim                = new(Certificate)
		key, crt, ca, crl, csr []byte
		block                  *pem.Block
	)

	for _, b := range inbound {
		for block, b = pem.Decode(b); block != nil; block, b = pem.Decode(b) {
			switch block.Type {
			case _CERTIFICATE:
				switch {
				case crt == nil:
					crt = block.Bytes
				default:
					ca = append(ca, block.Bytes...)
				}
			case _PRIVATE_KEY, _RSA_PRIVATE_KEY, _EC_PRIVATE_KEY:
				switch {
				case key != nil:
					return mod_errors.EParse
				}

				key = block.Bytes
			case _X509_CRL:
				switch {
				case crl != nil:
					return mod_errors.EParse
				}

				crl = block.Bytes
			case _CERTIFICATE_REQUEST:
				switch {
				case csr != nil:
					return mod_errors.EParse
				}

				csr = block.Bytes
			default:
				return mod_errors.EParse
			}
		}
	}

	switch err = interim.parseDER(key, crt, ca, csr, crl); {
	case err != nil:
		return
	}

	*r = *interim

	return
}

func (r *Certificate) parseP12(p12 []byte) (err error) {
	var (
		interim = new(Certificate)
	)
	switch interim.PrivateKey, interim.Certificate, interim.CertificateCAChain, err = pkcs12.DecodeChain(p12, pkcs12.DefaultPassword); {
	case err != nil:
		return
	}

	*r = *interim

	return
}

func (r *Certificate) encodeDER() (err error) {
	var (
		interim = new(Certificate)
	)

	switch {
	case r.PrivateKey != nil:
		switch interim.PrivateKeyDER, err = x509.MarshalPKCS8PrivateKey(r.PrivateKey); {
		case err != nil:
			return
		}
	}

	switch {
	case r.Certificate != nil:
		interim.CertificateDER = r.Certificate.Raw
	}

	switch {
	case r.CertificateCAChain != nil:
		interim.CertificateCAChainDER = make([][]byte, len(r.CertificateCAChain), len(r.CertificateCAChain))
		for a, b := range r.CertificateCAChain {
			interim.CertificateCAChainDER[a] = b.Raw
		}
	}

	switch {
	case r.CertificateRequest != nil:
		interim.CertificateRequestDER = r.CertificateRequest.Raw
	}

	switch {
	case r.RevocationList != nil:
		interim.RevocationListDER = r.RevocationList.Raw
	}

	r.PrivateKeyDER = interim.PrivateKeyDER
	r.CertificateDER = interim.CertificateDER
	r.CertificateCAChainDER = interim.CertificateCAChainDER
	r.CertificateRequestDER = interim.CertificateRequestDER
	r.RevocationListDER = interim.RevocationListDER

	return
}

func (r *Certificate) encodePEM() (err error) {
	var (
		interim = new(Certificate)
	)

	switch {
	case r.PrivateKey != nil:
		var (
			der []byte
		)
		switch der, err = x509.MarshalPKCS8PrivateKey(r.PrivateKey); {
		case err != nil:
			return
		}

		var (
			block = &pem.Block{
				Type:    _PRIVATE_KEY,
				Headers: nil,
				Bytes:   der,
			}
		)

		interim.PrivateKeyPEM = pem.EncodeToMemory(block)
	}

	switch {
	case r.Certificate != nil:
		var (
			block = &pem.Block{
				Type:    _CERTIFICATE,
				Headers: nil,
				Bytes:   r.Certificate.Raw,
			}
		)

		interim.CertificatePEM = pem.EncodeToMemory(block)
	}

	switch {
	case r.CertificateCAChain != nil:
		interim.CertificateCAChainPEM = make([][]byte, len(r.CertificateCAChain), len(r.CertificateCAChain))
		for a, b := range r.CertificateCAChain {
			var (
				block = &pem.Block{
					Type:    _CERTIFICATE,
					Headers: nil,
					Bytes:   b.Raw,
				}
			)

			interim.CertificateCAChainPEM[a] = pem.EncodeToMemory(block)
		}
	}

	switch {
	case r.CertificateRequest != nil:
		var (
			block = &pem.Block{
				Type:    _CERTIFICATE,
				Headers: nil,
				Bytes:   r.CertificateRequest.Raw,
			}
		)

		interim.CertificateRequestPEM = pem.EncodeToMemory(block)
	}

	switch {
	case r.RevocationList != nil:
		var (
			block = &pem.Block{
				Type:    _CERTIFICATE,
				Headers: nil,
				Bytes:   r.RevocationList.Raw,
			}
		)

		interim.RevocationListPEM = pem.EncodeToMemory(block)
	}

	r.PrivateKeyPEM = interim.PrivateKeyPEM
	r.CertificatePEM = interim.CertificatePEM
	r.CertificateCAChainPEM = interim.CertificateCAChainPEM
	r.CertificateRequestPEM = interim.CertificateRequestPEM
	r.RevocationListPEM = interim.RevocationListPEM

	r.PEM = bytes.Join([][]byte{
		r.PrivateKeyPEM,
		r.CertificatePEM,
		func() []byte { return bytes.Join(r.CertificateCAChainPEM, nil) }(),
		r.CertificateRequestPEM,
		r.RevocationListPEM,
	}, nil)

	return
}

func (r *Certificate) encodeP12() (err error) {
	var (
		interim []byte
	)
	switch interim, err = pkcs12.LegacyRC2.Encode(r.PrivateKey, r.Certificate, r.CertificateCAChain, pkcs12.DefaultPassword); {
	case err != nil:
		return
	}

	r.P12 = interim

	return
}
