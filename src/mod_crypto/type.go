package mod_crypto

import (
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type AuthDB map[string]transport.AuthMethod
