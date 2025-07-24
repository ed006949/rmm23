package mod_crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

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
	switch pub := r.Certificates[0].PublicKey.(type) {
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
