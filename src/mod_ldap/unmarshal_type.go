package mod_ldap

// LDAPAttributeUnmarshaler is the interface implemented by types
// that can UnmarshalEntry an LDAP attribute value representation of themselves.
type LDAPAttributeUnmarshaler interface {
	UnmarshalLDAPAttr([]string) error
}
