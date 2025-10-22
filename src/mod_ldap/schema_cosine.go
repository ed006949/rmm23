package mod_ldap

// COSINESchema returns the RFC 4524 COSINE LDAP schema.
func COSINESchema() *Schema {
	schema := NewSchema()

	// Register additional matching rules for COSINE
	registerCOSINEMatchingRules(schema)

	// Register COSINE attribute types
	registerCOSINEAttributeTypes(schema)

	// Register COSINE object classes
	registerCOSINEObjectClasses(schema)

	return schema
}

func registerCOSINEMatchingRules(s *Schema) {
	rules := []*MatchingRule{
		{OID: MatchingRuleOIDCaseIgnoreIA5Match, Names: []string{"caseIgnoreIA5Match"}, Syntax: SyntaxOIDIA5String},
		{OID: MatchingRuleOIDCaseIgnoreIA5SubstringsMatch, Names: []string{"caseIgnoreIA5SubstringsMatch"}, Syntax: SyntaxOIDIA5String},
	}

	for _, rule := range rules {
		s.MatchingRules[rule.OID] = rule
		for _, name := range rule.Names {
			s.MatchingRules[name] = rule
		}
	}
}

func registerCOSINEAttributeTypes(s *Schema) {
	attributes := []*AttributeType{
		{
			OID: AttributeOIDAssociatedDomain, Names: []string{"associatedDomain"},
			Equality: MatchingRuleOIDCaseIgnoreIA5Match, Substring: MatchingRuleOIDCaseIgnoreIA5SubstringsMatch,
			Syntax: SyntaxOIDIA5String,
		},
		{
			OID: AttributeOIDAssociatedName, Names: []string{"associatedName"},
			Equality: MatchingRuleOIDDistinguishedNameMatch, Syntax: SyntaxOIDDN,
		},
		{
			OID: AttributeOIDBuildingName, Names: []string{"buildingName"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDCo, Names: []string{"co", "friendlyCountryName"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDDocumentAuthor, Names: []string{"documentAuthor"},
			Equality: MatchingRuleOIDDistinguishedNameMatch, Syntax: SyntaxOIDDN,
		},
		{
			OID: AttributeOIDDocumentIdentifier, Names: []string{"documentIdentifier"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDDocumentLocation, Names: []string{"documentLocation"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDDocumentPublisher, Names: []string{"documentPublisher"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDDocumentTitle, Names: []string{"documentTitle"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDDocumentVersion, Names: []string{"documentVersion"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDDrink, Names: []string{"drink", "favouriteDrink"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDHomePhone, Names: []string{"homePhone", "homeTelephoneNumber"},
			Equality: MatchingRuleOIDTelephoneNumberMatch, Substring: MatchingRuleOIDTelephoneNumberSubstringsMatch,
			Syntax: SyntaxOIDTelephoneNumber,
		},
		{
			OID: AttributeOIDHomePostalAddress, Names: []string{"homePostalAddress"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDPostalAddress,
		},
		{
			OID: AttributeOIDHost, Names: []string{"host"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDInfo, Names: []string{"info"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDMail, Names: []string{"mail", "rfc822Mailbox"},
			Equality: MatchingRuleOIDCaseIgnoreIA5Match, Substring: MatchingRuleOIDCaseIgnoreIA5SubstringsMatch,
			Syntax: SyntaxOIDIA5String, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDManager, Names: []string{"manager"},
			Equality: MatchingRuleOIDDistinguishedNameMatch, Syntax: SyntaxOIDDN,
		},
		{
			OID: AttributeOIDMobile, Names: []string{"mobile", "mobileTelephoneNumber"},
			Equality: MatchingRuleOIDTelephoneNumberMatch, Substring: MatchingRuleOIDTelephoneNumberSubstringsMatch,
			Syntax: SyntaxOIDTelephoneNumber,
		},
		{
			OID: AttributeOIDOrganizationalStatus, Names: []string{"organizationalStatus"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDPager, Names: []string{"pager", "pagerTelephoneNumber"},
			Equality: MatchingRuleOIDTelephoneNumberMatch, Substring: MatchingRuleOIDTelephoneNumberSubstringsMatch,
			Syntax: SyntaxOIDTelephoneNumber,
		},
		{
			OID: AttributeOIDPersonalTitle, Names: []string{"personalTitle"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDRoomNumber, Names: []string{"roomNumber"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDSecretary, Names: []string{"secretary"},
			Equality: MatchingRuleOIDDistinguishedNameMatch, Syntax: SyntaxOIDDN,
		},
		{
			OID: AttributeOIDUniqueIdentifier, Names: []string{"uniqueIdentifier"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
		{
			OID: AttributeOIDUserClass, Names: []string{"userClass"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString, MaxLength: MaxLengthDirectoryString256,
		},
	}

	for _, attr := range attributes {
		s.AttributeTypes[attr.OID] = attr
		for _, name := range attr.Names {
			s.AttributeTypes[name] = attr
		}
	}
}

func registerCOSINEObjectClasses(s *Schema) {
	classes := []*ObjectClass{
		{
			OID: ObjectClassOIDAccount, Names: []string{"account"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"uid"},
			May:  []string{"description", "seeAlso", "l", "o", "ou", "host"},
		},
		{
			OID: ObjectClassOIDDocument, Names: []string{"document"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"documentIdentifier"},
			May: []string{"cn", "description", "seeAlso", "l", "o", "ou",
				"documentTitle", "documentVersion", "documentAuthor",
				"documentLocation", "documentPublisher"},
		},
		{
			OID: ObjectClassOIDDocumentSeries, Names: []string{"documentSeries"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"cn"},
			May:  []string{"description", "seeAlso", "telephonenumber", "l", "o", "ou"},
		},
		{
			OID: ObjectClassOIDDomain, Names: []string{"domain"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"dc"},
			May: []string{"userPassword", "searchGuide", "seeAlso", "businessCategory",
				"x121Address", "registeredAddress", "destinationIndicator",
				"preferredDeliveryMethod", "telexNumber", "teletexTerminalIdentifier",
				"telephoneNumber", "internationaliSDNNumber", "facsimileTelephoneNumber",
				"street", "postOfficeBox", "postalCode", "postalAddress",
				"physicalDeliveryOfficeName", "st", "l", "description", "o",
				"associatedName"},
		},
		{
			OID: ObjectClassOIDDomainRelatedObject, Names: []string{"domainRelatedObject"},
			Kind: ObjectClassKindAuxiliary, SuperClasses: []string{"top"},
			Must: []string{"associatedDomain"},
		},
		{
			OID: ObjectClassOIDFriendlyCountry, Names: []string{"friendlyCountry"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"country"},
			Must: []string{"co"},
		},
		{
			OID: ObjectClassOIDRFC822LocalPart, Names: []string{"rfc822LocalPart"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"domain"},
			Must: []string{"dc"},
			May: []string{"cn", "description", "destinationIndicator", "facsimileTelephoneNumber",
				"internationaliSDNNumber", "physicalDeliveryOfficeName", "postalAddress",
				"postalCode", "postOfficeBox", "preferredDeliveryMethod", "registeredAddress",
				"seeAlso", "sn", "street", "telephoneNumber", "teletexTerminalIdentifier",
				"telexNumber", "x121Address"},
		},
		{
			OID: ObjectClassOIDRoom, Names: []string{"room"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"cn"},
			May:  []string{"roomNumber", "description", "seeAlso", "telephoneNumber"},
		},
		{
			OID: ObjectClassOIDSimpleSecurityObject, Names: []string{"simpleSecurityObject"},
			Kind: ObjectClassKindAuxiliary, SuperClasses: []string{"top"},
			Must: []string{"userPassword"},
		},
	}

	for _, class := range classes {
		s.ObjectClasses[class.OID] = class
		for _, name := range class.Names {
			s.ObjectClasses[name] = class
		}
	}
}
