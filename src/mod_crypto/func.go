package mod_crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"

	"rmm23/src/mod_errors"
)

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

func parsePEM(inbound []byte) (outbound *Certificate, err error) {
	var (
		interim                = new(Certificate)
		key, crt, ca, crl, csr []byte
	)

	func() {
		var (
			certificateCounter int
		)

		for len(inbound) != 0 {
			var (
				interimDERBlock *pem.Block
			)

			interimDERBlock, inbound = pem.Decode(inbound)
			switch {
			case interimDERBlock == nil:
				return
			}

			switch {
			case interimDERBlock.Type == _CERTIFICATE:
				switch certificateCounter {
				case 0:
					crt = interimDERBlock.Bytes
				default:
					ca = append(ca, interimDERBlock.Bytes...)
				}

				certificateCounter++
			case interimDERBlock.Type == _PRIVATE_KEY || strings.HasSuffix(interimDERBlock.Type, __PRIVATE_KEY):
				key = interimDERBlock.Bytes
			}
		}
	}()

	switch err = interim.ParseDERs(key, crt, ca, crl, csr); {
	case err == nil:
		return
	}

	return interim, nil
}
