package mod_crypto

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"

	"rmm23/src/mod_errors"
)

//
// taken from https://github.com/golang/go/src/crypto/tls/tls.go
// modified to be more useful
//

func X509KeyPair(certPEMBlock []byte, keyPEMBlock []byte) (outbound *Certificate, err error) {
	var (
		certDERBlock *pem.Block
		keyDERBlock  *pem.Block
		interimCert  *x509.Certificate
	)
	outbound = new(Certificate)

	func() {
		for {
			certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
			// for ; certDERBlock != nil; certDERBlock, certPEMBlock = pem.Decode(certPEMBlock) {
			switch {
			case certDERBlock == nil:
				return
			case certDERBlock.Type == "CERTIFICATE":
				outbound.CertificatesDER = append(
					outbound.CertificatesDER,
					certDERBlock.Bytes,
				)
				outbound.CertificatesRawPEM = append(
					outbound.CertificatesRawPEM,
					[]byte(base64.RawStdEncoding.EncodeToString(certDERBlock.Bytes)),
				)

				switch {
				case len(outbound.CertificatesDER) > 1:
					outbound.CertificateCAChainDER = append(outbound.CertificateCAChainDER, certDERBlock.Bytes...)
				}
			}
		}
	}()
	switch {
	case len(outbound.CertificatesDER) == 0:
		return nil, mod_errors.EPEMNoDataCert
	}
	outbound.CertificateCAChainRawPEM = []byte(base64.RawStdEncoding.EncodeToString(outbound.CertificateCAChainDER))

	func() {
		for {
			keyDERBlock, keyPEMBlock = pem.Decode(keyPEMBlock)
			// for ; keyDERBlock != nil; keyDERBlock, keyPEMBlock = pem.Decode(keyPEMBlock) {
			switch {
			case keyDERBlock == nil:
				return
			case keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY"):
				outbound.PrivateKeyDER = keyDERBlock.Bytes
				outbound.PrivateKeyRawPEM = []byte(base64.RawStdEncoding.EncodeToString(outbound.PrivateKeyDER))
			}
		}
	}()
	switch {
	case len(outbound.PrivateKeyDER) == 0:
		return nil, mod_errors.EPEMNoDataKey
	}

	for _, b := range outbound.CertificatesDER {
		switch interimCert, err = x509.ParseCertificate(b); {
		case err != nil:
			return nil, err
		default:
			outbound.Certificates = append(outbound.Certificates, interimCert)
		}
	}

	switch outbound.PrivateKey, err = ParsePrivateKey(outbound.PrivateKeyDER); {
	case err != nil:
		return nil, err
	}

	// TODO complete local chain verification
	// We don't need to parse the public key for TLS, but we so do anyway
	// to check that it looks sane and matches the private key.

	switch pub := outbound.Certificates[0].PublicKey.(type) {
	case *rsa.PublicKey:
		switch priv, ok := outbound.PrivateKey.(*rsa.PrivateKey); {
		case !ok:
			return nil, mod_errors.ETypeMismatchPrivKeyPubKey
		case pub.N.Cmp(priv.N) != 0:
			return nil, mod_errors.EMismatchPrivKeyPubKey
		}
	case *ecdsa.PublicKey:
		switch priv, ok := outbound.PrivateKey.(*ecdsa.PrivateKey); {
		case !ok:
			return nil, mod_errors.ETypeMismatchPrivKeyPubKey
		case pub.X.Cmp(priv.X) != 0 || pub.Y.Cmp(priv.Y) != 0:
			return nil, mod_errors.EMismatchPrivKeyPubKey
		}
	case ed25519.PublicKey:
		switch priv, ok := outbound.PrivateKey.(ed25519.PrivateKey); {
		case !ok:
			return nil, mod_errors.ETypeMismatchPrivKeyPubKey
		case !bytes.Equal(priv.Public().(ed25519.PublicKey), pub):
			return nil, mod_errors.EMismatchPrivKeyPubKey
		}
	default:
		return nil, mod_errors.EUnknownPubKeyAlgo
	}

	return
}
func ParsePrivateKey(der []byte) (key crypto.PrivateKey, err error) {
	switch key, err = x509.ParsePKCS1PrivateKey(der); {
	case err == nil:
		return
	}

	switch key, err = x509.ParsePKCS8PrivateKey(der); {
	case err == nil:
		switch value := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
			return value, nil
		default:
			return nil, mod_errors.EUnknownPrivKeyType
		}
	}

	switch key, err = x509.ParseECPrivateKey(der); {
	case err == nil:
		return
	}

	return nil, mod_errors.EX509ParsePrivKey
}

func ParsePEM(PEMBlock []byte) (outbound *Certificate, err error) {
	var (
		interimCert *x509.Certificate
	)
	outbound = new(Certificate)

	func() {
		for {
			var (
				interimDERBlock *pem.Block
			)
			interimDERBlock, PEMBlock = pem.Decode(PEMBlock)
			switch {
			case interimDERBlock == nil:
				return
			case interimDERBlock.Type == "CERTIFICATE":
				outbound.CertificatesDER = append(
					outbound.CertificatesDER,
					interimDERBlock.Bytes,
				)
				outbound.CertificatesRawPEM = append(
					outbound.CertificatesRawPEM,
					[]byte(base64.RawStdEncoding.EncodeToString(interimDERBlock.Bytes)),
				)

				switch {
				case len(outbound.CertificatesDER) > 1:
					outbound.CertificateCAChainDER = append(outbound.CertificateCAChainDER, interimDERBlock.Bytes...)
				}
			case interimDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(interimDERBlock.Type, " PRIVATE KEY"):
				outbound.PrivateKeyDER = interimDERBlock.Bytes
				outbound.PrivateKeyRawPEM = []byte(base64.RawStdEncoding.EncodeToString(outbound.PrivateKeyDER))
			}
		}
	}()
	switch {
	case len(outbound.CertificatesDER) == 0:
		return nil, mod_errors.EPEMNoDataCert
	case len(outbound.PrivateKeyDER) == 0:
		return nil, mod_errors.EPEMNoDataKey
	}

	outbound.CertificateCAChainRawPEM = []byte(base64.RawStdEncoding.EncodeToString(outbound.CertificateCAChainDER))

	for _, b := range outbound.CertificatesDER {
		switch interimCert, err = x509.ParseCertificate(b); {
		case err != nil:
			return nil, err
		default:
			outbound.Certificates = append(outbound.Certificates, interimCert)
		}
	}

	switch outbound.PrivateKey, err = ParsePrivateKey(outbound.PrivateKeyDER); {
	case err != nil:
		return nil, err
	}

	// TODO complete local chain verification
	// We don't need to parse the public key for TLS, but we so do anyway
	// to check that it looks sane and matches the private key.

	switch pub := outbound.Certificates[0].PublicKey.(type) {
	case *rsa.PublicKey:
		switch priv, ok := outbound.PrivateKey.(*rsa.PrivateKey); {
		case !ok:
			return nil, mod_errors.ETypeMismatchPrivKeyPubKey
		case pub.N.Cmp(priv.N) != 0:
			return nil, mod_errors.EMismatchPrivKeyPubKey
		}
	case *ecdsa.PublicKey:
		switch priv, ok := outbound.PrivateKey.(*ecdsa.PrivateKey); {
		case !ok:
			return nil, mod_errors.ETypeMismatchPrivKeyPubKey
		case pub.X.Cmp(priv.X) != 0 || pub.Y.Cmp(priv.Y) != 0:
			return nil, mod_errors.EMismatchPrivKeyPubKey
		}
	case ed25519.PublicKey:
		switch priv, ok := outbound.PrivateKey.(ed25519.PrivateKey); {
		case !ok:
			return nil, mod_errors.ETypeMismatchPrivKeyPubKey
		case !bytes.Equal(priv.Public().(ed25519.PublicKey), pub):
			return nil, mod_errors.EMismatchPrivKeyPubKey
		}
	default:
		return nil, mod_errors.EUnknownPubKeyAlgo
	}

	return
}
