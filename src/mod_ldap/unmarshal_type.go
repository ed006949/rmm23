package mod_ldap

import (
	"encoding/xml"
)

// LDAPAttributeUnmarshaler is the interface implemented by types
// that can unmarshal an LDAP attribute value representation of themselves.
type LDAPAttributeUnmarshaler interface {
	UnmarshalLDAPAttr([]string) error
}

type GenericXMLElement struct {
	XMLName    xml.Name
	Attrs      []xml.Attr          `xml:",any,attr"`
	InnerText  string              `xml:",chardata"`
	ChildNodes []GenericXMLElement `xml:",any"`
}
