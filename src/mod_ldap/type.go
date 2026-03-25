package mod_ldap

import (
	"github.com/go-ldap/ldap/v3"

	"rmm23/src/mod_url"
)

type Conf struct {
	URL      *mod_url.URL `json:"url,omitempty"`
	Domains  []*domain    `json:"domain,omitempty"`
	Settings []*settings  `json:"settings,omitempty"`
	conn     *ldap.Conn
}

type domain struct {
	DN string `json:"dn,omitempty"`
	// Settings []*settings `json:"settings,omitempty"`
}

type settings struct {
	Type   string           `json:"type,omitempty"`
	DN     string           `json:"dn,omitempty"`
	CN     string           `json:"cn,omitempty"`
	Scope  attrSearchScope  `json:"scope,omitempty"`
	Filter attrSearchFilter `json:"filter,omitempty"`
}

type attrSearchScope int
type attrSearchFilter string
