package mod_ldap

// OpenSSHLPKSchema returns the OpenSSH LDAP Public Key schema
// This schema allows storing SSH public keys in LDAP for SSH authentication.
func OpenSSHLPKSchema() *Schema {
	schema := NewSchema()

	// Register OpenSSH LPK attribute types
	registerOpenSSHLPKAttributeTypes(schema)

	// Register OpenSSH LPK object classes
	registerOpenSSHLPKObjectClasses(schema)

	return schema
}

func registerOpenSSHLPKAttributeTypes(s *Schema) {
	attributes := []*AttributeType{
		{
			OID:         AttributeOIDSshPublicKey,
			Names:       []string{"sshPublicKey"},
			Description: "OpenSSH public key",
			Equality:    MatchingRuleOIDOctetStringMatch,
			Syntax:      SyntaxOIDOctetString,
		},
	}

	for _, attr := range attributes {
		s.AttributeTypes[attr.OID] = attr
		for _, name := range attr.Names {
			s.AttributeTypes[name] = attr
		}
	}
}

func registerOpenSSHLPKObjectClasses(s *Schema) {
	classes := []*ObjectClass{
		{
			OID:          ObjectClassOIDLdapPublicKey,
			Names:        []string{"ldapPublicKey"},
			Description:  "OpenSSH LDAP Public Key",
			Kind:         ObjectClassKindAuxiliary,
			SuperClasses: []string{"top"},
			Must:         []string{"sshPublicKey"},
			May:          []string{"uid"},
		},
	}

	for _, class := range classes {
		s.ObjectClasses[class.OID] = class
		for _, name := range class.Names {
			s.ObjectClasses[name] = class
		}
	}
}
