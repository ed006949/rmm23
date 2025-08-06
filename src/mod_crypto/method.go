package mod_crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"software.sslmate.com/src/go-pkcs12"

	"rmm23/src/mod_errors"
)

func (r *AuthDB) WriteSSH(name string, user string, pemBytes []byte, password string) (err error) {
	switch _, ok := (*r)[name]; {
	case ok:
		return mod_errors.EDUPDATA
	}

	var (
		sshPublicKeys *ssh.PublicKeys
	)
	switch sshPublicKeys, err = ssh.NewPublicKeys(user, pemBytes, password); {
	case err != nil:
		return
	default:
		(*r)[name] = sshPublicKeys

		return
	}
}

func (r *AuthDB) WriteToken(name string, user string, tokenBytes []byte) (err error) {
	switch _, ok := (*r)[name]; {
	case ok:
		return mod_errors.EDUPDATA
	}

	(*r)[name] = &http.BasicAuth{
		Username: user,
		Password: string(tokenBytes),
	}

	return
}

func (r *AuthDB) ReadAuth(name string) (outbound transport.AuthMethod, err error) {
	switch value, ok := (*r)[name]; {
	case !ok:
		return nil, mod_errors.ENOTFOUND
	default:
		return value, nil
	}
}

// checkPrivateKey
// We don't need to parse the public key for TLS, but we so do anyway to check that it looks sane and matches the private key.
func (r *Certificate) checkPrivateKey() (err error) {
	switch pub := r.Certificate.PublicKey.(type) {
	case *rsa.PublicKey:
		return r.checkPrivateKeyRSA(pub)
	case *ecdsa.PublicKey:
		return r.checkPrivateKeyECDSA(pub)
	case ed25519.PublicKey:
		return r.checkPrivateKeyED25519(pub)
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

func (r *Certificate) decodeP12() (err error) {
	var (
		privateKey  any
		certificate *x509.Certificate
		chain       []*x509.Certificate
	)
	switch privateKey, certificate, chain, err = pkcs12.DecodeChain(r.P12, pkcs12.DefaultPassword); {
	case err != nil:
		return
	}

	r.PrivateKey = privateKey
	r.Certificate = certificate
	r.CertificateCAChain = chain

	return
}

func (r *Certificate) ParseDERs(key, crt, ca, crl, csr []byte) (err error) {
	var (
		interim = new(Certificate)
	)

	switch interim.PrivateKey, err = x509.ParsePKCS8PrivateKey(key); {
	case err != nil:
		return
	}

	switch interim.Certificate, err = x509.ParseCertificate(crt); {
	case err != nil:
		return
	}

	switch interim.CertificateCAChain, err = x509.ParseCertificates(ca); {
	case err != nil:
		// return
	default:
		for _, b := range interim.CertificateCAChain {
			interim.CertificateCAChainDER = append(interim.CertificateCAChainDER, b.Raw)
		}
	}

	switch interim.RevocationList, err = x509.ParseRevocationList(crl); {
	case err != nil:
	// return
	default:
		interim.RevocationListDER = interim.RevocationList.Raw
	}

	switch interim.CertificateRequest, err = x509.ParseCertificateRequest(csr); {
	case err != nil:
		// return
	default:
		interim.CertificateRequestDER = interim.CertificateRequest.Raw
	}

	switch err = interim.checkPrivateKey(); {
	case err != nil:
		return
	}

	switch err = interim.encodeP12(); {
	case err != nil:
		return
	}

	*r = *interim

	return
}

func (r *Certificate) ParsePEM(inbound []byte) (err error) {
	var (
		interim                = new(Certificate)
		key, crt, ca, crl, csr []byte
		block                  *pem.Block
	)

	for block, inbound = pem.Decode(inbound); block != nil; block, inbound = pem.Decode(inbound) {
		switch block.Type {
		case _CERTIFICATE:
			switch {
			case crt == nil:
				crt = block.Bytes
			default:
				ca = append(ca, block.Bytes...)
			}
		case _PRIVATE_KEY, _DSA_PRIVATE_KEY, _RSA_PRIVATE_KEY, _EC_PRIVATE_KEY:
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

	switch err = interim.ParseDERs(key, crt, ca, crl, csr); {
	case err == nil:
		return
	}

	*r = *interim

	return
}
