package mod_ldap

import (
	"github.com/go-ldap/ldap/v3"

	"rmm23/src/mod_net"
)

// LDAPAttributeUnmarshaler is the interface implemented by types
// that can UnmarshalEntry an LDAP attribute value representation of themselves.
type LDAPAttributeUnmarshaler interface {
	UnmarshalLDAPAttr(values []string) (err error)
}

type Conf struct {
	URL      *mod_net.URL `json:"url,omitempty"`
	Settings []*settings  `json:"settings,omitempty"`
	Domains  []*domain    `json:"domain,omitempty"`
	conn     *ldap.Conn
}

type settings struct {
	Type   string           `json:"type,omitempty"`
	DN     string           `json:"dn,omitempty"`
	CN     string           `json:"cn,omitempty"`
	Scope  attrSearchScope  `json:"scope,omitempty"`
	Filter attrSearchFilter `json:"filter,omitempty"`
}

type domain struct {
	DN string `json:"dn,omitempty"`

	// Settings []*settings `json:"settings,omitempty"`
}

type attrSearchScope int
type attrSearchFilter string
