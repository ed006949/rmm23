package mod_ldap

// LabeledURISchema returns the RFC 2079 labeledURI attribute schema.
func LabeledURISchema() *Schema {
	schema := NewSchema()

	// Register labeledURI attribute type
	registerLabeledURIAttributeTypes(schema)

	return schema
}

func registerLabeledURIAttributeTypes(s *Schema) {
	attributes := []*AttributeType{
		{
			OID:         AttributeOIDLabeledURI,
			Names:       []string{"labeledURI"},
			Description: "Uniform Resource Identifier with optional label",
			Equality:    MatchingRuleOIDCaseExactMatch,
			Syntax:      SyntaxOIDDirectoryString,
		},
	}

	for _, attr := range attributes {
		s.AttributeTypes[attr.OID] = attr
		for _, name := range attr.Names {
			s.AttributeTypes[name] = attr
		}
	}
}
