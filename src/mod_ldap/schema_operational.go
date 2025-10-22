package mod_ldap

// OperationalSchema returns the operational attributes schema
// These attributes are automatically maintained by LDAP servers.
func OperationalSchema() *Schema {
	schema := NewSchema()

	// Register operational attribute types
	registerOperationalAttributeTypes(schema)

	return schema
}

func registerOperationalAttributeTypes(s *Schema) {
	attributes := []*AttributeType{
		// Standard operational attributes (RFC 4512)
		{
			OID:          AttributeOIDCreateTimestamp,
			Names:        []string{"createTimestamp"},
			Description:  "Time when the entry was added to the directory",
			Equality:     MatchingRuleOIDGeneralizedTimeMatch,
			Ordering:     MatchingRuleOIDGeneralizedTimeOrderingMatch,
			Syntax:       SyntaxOIDGeneralizedTime,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},
		{
			OID:          AttributeOIDModifyTimestamp,
			Names:        []string{"modifyTimestamp"},
			Description:  "Time when the entry was last modified",
			Equality:     MatchingRuleOIDGeneralizedTimeMatch,
			Ordering:     MatchingRuleOIDGeneralizedTimeOrderingMatch,
			Syntax:       SyntaxOIDGeneralizedTime,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},
		{
			OID:          AttributeOIDCreatorsName,
			Names:        []string{"creatorsName"},
			Description:  "DN of the user who added this entry to the directory",
			Equality:     MatchingRuleOIDDistinguishedNameMatch,
			Syntax:       SyntaxOIDDN,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},
		{
			OID:          AttributeOIDModifiersName,
			Names:        []string{"modifiersName"},
			Description:  "DN of the user who last modified this entry",
			Equality:     MatchingRuleOIDDistinguishedNameMatch,
			Syntax:       SyntaxOIDDN,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},
		{
			OID:          AttributeOIDSubschemaSubentry,
			Names:        []string{"subschemaSubentry"},
			Description:  "DN of the subschema entry controlling this entry",
			Equality:     MatchingRuleOIDDistinguishedNameMatch,
			Syntax:       SyntaxOIDDN,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},
		{
			OID:          AttributeOIDStructuralObjectClass,
			Names:        []string{"structuralObjectClass"},
			Description:  "Structural object class of the entry",
			Equality:     MatchingRuleOIDObjectIdentifierMatch,
			Syntax:       SyntaxOIDOID,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},

		// entryUUID (RFC 4530)
		{
			OID:          AttributeOIDEntryUUID,
			Names:        []string{"entryUUID"},
			Description:  "UUID of the entry",
			Equality:     MatchingRuleOIDUniqueIdentifierMatch,
			Ordering:     MatchingRuleOIDUniqueIdentifierMatch,
			Syntax:       SyntaxOIDBitString,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},

		// Common operational attributes (widely supported)
		{
			OID:          AttributeOIDEntryDN,
			Names:        []string{"entryDN"},
			Description:  "DN of the entry",
			Equality:     MatchingRuleOIDDistinguishedNameMatch,
			Syntax:       SyntaxOIDDN,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},
		{
			OID:          AttributeOIDHasSubordinates,
			Names:        []string{"hasSubordinates"},
			Description:  "Whether this entry has subordinate entries",
			Equality:     MatchingRuleOIDBooleanMatch,
			Syntax:       SyntaxOIDBoolean,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},
		{
			OID:          AttributeOIDSubordinateCount,
			Names:        []string{"subordinateCount"},
			Description:  "Number of subordinate entries",
			Equality:     MatchingRuleOIDIntegerMatch,
			Ordering:     MatchingRuleOIDIntegerOrderingMatch,
			Syntax:       SyntaxOIDInteger,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},
		{
			OID:          AttributeOIDEntryCSN,
			Names:        []string{"entryCSN"},
			Description:  "Change sequence number of the entry (for replication)",
			Equality:     MatchingRuleOIDOctetStringMatch,
			Ordering:     MatchingRuleOIDOctetStringOrderingMatch,
			Syntax:       SyntaxOIDOctetString,
			SingleValue:  true,
			NoUserModify: true,
			Usage:        AttributeUsageDirectoryOperation,
		},
		{
			OID:          AttributeOIDContextCSN,
			Names:        []string{"contextCSN"},
			Description:  "Context change sequence number (for replication)",
			Equality:     MatchingRuleOIDOctetStringMatch,
			Ordering:     MatchingRuleOIDOctetStringOrderingMatch,
			Syntax:       SyntaxOIDOctetString,
			NoUserModify: true,
			Usage:        AttributeUsageDSAOperation,
		},
	}

	for _, attr := range attributes {
		s.AttributeTypes[attr.OID] = attr
		for _, name := range attr.Names {
			s.AttributeTypes[name] = attr
		}
	}
}
