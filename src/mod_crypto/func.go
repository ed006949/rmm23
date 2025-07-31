package mod_crypto

import (
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

func parsePrivateKey(der []byte) (key crypto.PrivateKey, err error) {
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

func parsePEM(inbound []byte) (outbound *Certificate, err error) {
	var (
		interim = new(Certificate)
	)

	func() {
		var (
			certificateCounter int
		)

		for {
			var (
				interimDERBlock    *pem.Block
				interimCertificate *x509.Certificate
			)

			interimDERBlock, inbound = pem.Decode(inbound)

			switch {
			case interimDERBlock == nil:
				return
			}

			var (
				interimPEMBlock = []byte(base64.RawStdEncoding.EncodeToString(interimDERBlock.Bytes)) // sanitized PEM
			)

			switch {
			case interimDERBlock.Type == _CERTIFICATE:
				switch certificateCounter++; certificateCounter {
				case 1:
					switch interim.Certificate, err = x509.ParseCertificate(interimDERBlock.Bytes); {
					case err != nil:
						return
					}

					interim.CertificateDER = interimDERBlock.Bytes
					interim.CertificatePEM = interimPEMBlock
					interim.PEM = append(interim.PEM, interimPEMBlock...)
				default:
					switch interimCertificate, err = x509.ParseCertificate(interimDERBlock.Bytes); {
					case err != nil:
						return
					}

					interim.CertificateCAChain = append(interim.CertificateCAChain, interimCertificate)
					interim.CertificateCAChainDER = append(interim.CertificateCAChainDER, interimDERBlock.Bytes)
					interim.CertificateCAChainPEM = append(interim.CertificateCAChainPEM, interimPEMBlock)
					interim.PEM = append(interim.PEM, interimPEMBlock...)
				}
			case interimDERBlock.Type == _PRIVATE_KEY || strings.HasSuffix(interimDERBlock.Type, __PRIVATE_KEY):
				switch interim.PrivateKey, err = parsePrivateKey(interimDERBlock.Bytes); {
				case err != nil:
					return
				}

				interim.PrivateKeyDER = interimDERBlock.Bytes
				interim.PrivateKeyPEM = interimPEMBlock
				interim.PEM = append(interim.PEM, interimPEMBlock...)
			}
		}
	}()

	switch {
	case interim.Certificate == nil:
		return nil, mod_errors.EPEMNoDataCert
	case interim.PrivateKey == nil:
		return nil, mod_errors.EPEMNoDataKey
	}

	switch err = interim.checkPrivateKey(); {
	case err != nil:
		return nil, err
	}

	switch err = interim.EncodeP12(); {
	case err != nil:
		return nil, err
	}

	return interim, err
}
