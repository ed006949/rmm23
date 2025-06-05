package mod_crypto

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"rmm23/src/l"
)

func (r *AuthDB) WriteSSH(name string, user string, pemBytes []byte, password string) (err error) {
	switch _, ok := (*r)[name]; {
	case ok:
		return l.EDUPDATA
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
		return l.EDUPDATA
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
		return nil, l.ENOTFOUND
	default:
		return value, nil
	}
}
