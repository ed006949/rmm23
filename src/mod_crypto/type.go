package mod_crypto

import (
	"crypto"
	"crypto/x509"

	"github.com/go-git/go-git/v5/plumbing/transport"
)

type AuthDB map[string]transport.AuthMethod

type PemCertKey string
type PemCertKeyList map[string]*PemCertKey
type SignatureScheme uint16

type Certificate struct {
	PrivateKeyDER   []byte   `redis:"priv" redisearch:"text,sortable"`
	CertificatesDER [][]byte `redis:"crt" redisearch:"text,sortable"`
	// CertificateCAChainPEM    []byte

	PrivateKeyPEM         []byte
	CertificatesPEM       [][]byte
	CertificateCAChainDER []byte

	PrivateKeyRawPEM         []byte
	CertificatesRawPEM       [][]byte
	CertificateCAChainRawPEM []byte

	// Certificates is the parsed form of the leaf certificate, which may be initialized
	// using x509.ParseCertificate to reduce per-handshake processing. If nil,
	// the leaf certificate will be parsed as needed.
	Certificates []*x509.Certificate
	// PrivateKey contains the private key corresponding to the public key in
	// Certificates. This must implement crypto.Signer with an RSA, ECDSA or Ed25519 PublicKey.
	// For a server up to TLS 1.2, it can also implement crypto.Decrypter with
	// an RSA PublicKey.
	PrivateKey crypto.PrivateKey

	// SupportedSignatureAlgorithms is an optional list restricting what
	// signature algorithms the PrivateKey can be used for.
	SupportedSignatureAlgorithms []SignatureScheme
	// OCSPStaple contains an optional OCSP response which will be served
	// to clients that request it.
	OCSPStaple []byte
	// SignedCertificateTimestamps contains an optional list of Signed
	// Certificate Timestamps which will be served to clients that request it.
	SignedCertificateTimestamps [][]byte
}
