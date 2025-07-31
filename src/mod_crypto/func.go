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

//
// taken from https://github.com/golang/go/src/crypto/tls/tls.go
// modified to be more useful
//

// func X509KeyPair(certPEMBlock []byte, keyPEMBlock []byte) (outbound *Certificate, err error) {
// 	var (
// 		certDERBlock *pem.Block
// 		keyDERBlock  *pem.Block
// 		interimCert  *x509.Certificate
// 	)
//
// 	outbound = new(Certificate)
//
// 	func() {
// 		for {
// 			certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
// 			// for ; certDERBlock != nil; certDERBlock, certPEMBlock = pem.Decode(certPEMBlock) {
// 			switch {
// 			case certDERBlock == nil:
// 				return
// 			case certDERBlock.Type == _CERTIFICATE:
// 				outbound.CertificateDER = append(
// 					outbound.CertificateDER,
// 					certDERBlock.Bytes,
// 				)
// 				outbound.CertificateRawPEM = append(
// 					outbound.CertificateRawPEM,
// 					[]byte(base64.RawStdEncoding.EncodeToString(certDERBlock.Bytes)),
// 				)
//
// 				switch {
// 				case len(outbound.CertificateDER) > 1:
// 					outbound.CertificateCAChainDER = append(outbound.CertificateCAChainDER, certDERBlock.Bytes...)
// 				}
// 			}
// 		}
// 	}()
//
// 	switch {
// 	case len(outbound.CertificateDER) == 0:
// 		return nil, mod_errors.EPEMNoDataCert
// 	}
//
// 	outbound.CertificateCAChainPEM = append(outbound.CertificateCAChainPEM, []byte(base64.RawStdEncoding.EncodeToString(outbound.CertificateCAChainDER)))
//
// 	func() {
// 		for {
// 			keyDERBlock, keyPEMBlock = pem.Decode(keyPEMBlock)
// 			// for ; keyDERBlock != nil; keyDERBlock, keyPEMBlock = pem.Decode(keyPEMBlock) {
// 			switch {
// 			case keyDERBlock == nil:
// 				return
// 			case keyDERBlock.Type == _PRIVATE_KEY || strings.HasSuffix(keyDERBlock.Type, __PRIVATE_KEY):
// 				outbound.PrivateKeyDER = keyDERBlock.Bytes
// 				outbound.PrivateKeyRawPEM = []byte(base64.RawStdEncoding.EncodeToString(outbound.PrivateKeyDER))
// 			}
// 		}
// 	}()
//
// 	switch {
// 	case len(outbound.PrivateKeyDER) == 0:
// 		return nil, mod_errors.EPEMNoDataKey
// 	}
//
// 	switch outbound.Certificate, err = x509.ParseCertificate(outbound.CertificateDER); {
// 	case err != nil:
// 		return nil, err
// 	}
//
// 	for _, b := range outbound.CertificateCAChainDER {
// 		switch interimCert, err = x509.ParseCertificate(b); {
// 		case err != nil:
// 			return nil, err
// 		}
//
// 		outbound.CertificateCAChain = append(outbound.CertificateCAChain, interimCert)
// 	}
//
// 	switch outbound.PrivateKey, err = parsePrivateKey(outbound.PrivateKeyDER); {
// 	case err != nil:
// 		return nil, err
// 	}
//
// 	switch err = outbound.checkPrivateKey(); {
// 	case err != nil:
// 		return nil, err
// 	}
//
// 	return
// }

func parsePrivateKey(der []byte) (key crypto.PrivateKey, err error) {
	switch key, err = x509.ParsePKCS1PrivateKey(der); err {
	case nil:
		return
	}

	switch key, err = x509.ParsePKCS8PrivateKey(der); err {
	case nil:
		switch value := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
			return value, nil
		default:
			return nil, mod_errors.EUnknownPrivKeyType
		}
	}

	switch key, err = x509.ParseECPrivateKey(der); err {
	case nil:
		return
	}

	return nil, mod_errors.EX509ParsePrivKey
}

func parsePEM(inbound []byte) (outbound *Certificate, err error) {
	var (
		interim = new(Certificate)
	)

	func() {
		for {
			var (
				interimDERBlock    *pem.Block
				interimDERBlocks   []*pem.Block
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
				interimDERBlocks = append(interimDERBlocks, interimDERBlock)

				switch len(interimDERBlocks) {
				case 0:
					return
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
