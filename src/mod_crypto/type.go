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

type Certificates map[string]*Certificate

type Certificate struct {
	P12 []byte `json:"p12,omitempty"`
	DER []byte `json:"-"`
	PEM []byte `json:"-"`
	CRL []byte `json:"crl,omitempty"`
	CSR []byte `json:"csr,omitempty"`

	PrivateKeyDER         []byte   `json:"-"`
	CertificateRequestDER []byte   `json:"-"`
	CertificateDER        []byte   `json:"-"`
	CertificateCAChainDER [][]byte `json:"-"`
	RevocationListDER     []byte   `json:"-"`

	PrivateKeyPEM         []byte   `json:"-"`
	CertificateRequestPEM []byte   `json:"-"`
	CertificatePEM        []byte   `json:"-"`
	CertificateCAChainPEM [][]byte `json:"-"`
	RevocationListPEM     []byte   `json:"-"`

	// PrivateKey contains the private key corresponding to the public key in
	// Certificates. This must implement crypto.Signer with an RSA, ECDSA or Ed25519 PublicKey.
	// For a server up to TLS 1.2, it can also implement crypto.Decrypter with
	// an RSA PublicKey.
	PrivateKey crypto.PrivateKey `json:"-"`
	// Certificates is the parsed form of the leaf certificate, which may be initialized
	// using x509.ParseCertificate to reduce per-handshake processing. If nil,
	// the leaf certificate will be parsed as needed.
	// Certificates       []*x509.Certificate `json:"-"`
	CertificateRequest *x509.CertificateRequest `json:"-"`
	Certificate        *x509.Certificate        `json:"-"`
	CertificateCAChain []*x509.Certificate      `json:"-"`
	RevocationList     *x509.RevocationList     `json:"-"`
}
