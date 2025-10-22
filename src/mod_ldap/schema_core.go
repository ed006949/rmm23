package mod_ldap

// CoreSchema returns the RFC 4519 core LDAP schema.
func CoreSchema() *Schema {
	schema := NewSchema()

	// Register syntaxes
	registerCoreSyntaxes(schema)

	// Register matching rules
	registerCoreMatchingRules(schema)

	// Register attribute types
	registerCoreAttributeTypes(schema)

	// Register object classes
	registerCoreObjectClasses(schema)

	return schema
}

func registerCoreSyntaxes(s *Schema) {
	syntaxes := []*Syntax{
		{OID: SyntaxOIDAttributeTypeDescription, Description: "Attribute Type Description"},
		{OID: SyntaxOIDBitString, Description: "Bit String"},
		{OID: SyntaxOIDBoolean, Description: "Boolean"},
		{OID: SyntaxOIDCountryString, Description: "Country String"},
		{OID: SyntaxOIDDN, Description: "DN"},
		{OID: SyntaxOIDDirectoryString, Description: "Directory String"},
		{OID: SyntaxOIDDITContentRuleDescription, Description: "DIT Content Rule Description"},
		{OID: SyntaxOIDDITStructureRuleDescription, Description: "DIT Structure Rule Description"},
		{OID: SyntaxOIDEnhancedGuide, Description: "Enhanced Guide"},
		{OID: SyntaxOIDFacsimileTelephoneNumber, Description: "Facsimile Telephone Number"},
		{OID: SyntaxOIDFax, Description: "Fax"},
		{OID: SyntaxOIDGeneralizedTime, Description: "Generalized Time"},
		{OID: SyntaxOIDGuide, Description: "Guide"},
		{OID: SyntaxOIDIA5String, Description: "IA5 String"},
		{OID: SyntaxOIDInteger, Description: "Integer"},
		{OID: SyntaxOIDJPEG, Description: "JPEG"},
		{OID: SyntaxOIDMatchingRuleDescription, Description: "Matching Rule Description"},
		{OID: SyntaxOIDMatchingRuleUseDescription, Description: "Matching Rule Use Description"},
		{OID: SyntaxOIDNameAndOptionalUID, Description: "Name And Optional UID"},
		{OID: SyntaxOIDNameFormDescription, Description: "Name Form Description"},
		{OID: SyntaxOIDNumericString, Description: "Numeric String"},
		{OID: SyntaxOIDObjectClassDescription, Description: "Object Class Description"},
		{OID: SyntaxOIDOctetString, Description: "Octet String"},
		{OID: SyntaxOIDOID, Description: "OID"},
		{OID: SyntaxOIDPostalAddress, Description: "Postal Address"},
		{OID: SyntaxOIDPrintableString, Description: "Printable String"},
		{OID: SyntaxOIDTelephoneNumber, Description: "Telephone Number"},
		{OID: SyntaxOIDTeletexTerminalIdentifier, Description: "Teletex Terminal Identifier"},
		{OID: SyntaxOIDTelexNumber, Description: "Telex Number"},
	}

	for _, syntax := range syntaxes {
		s.Syntaxes[syntax.OID] = syntax
	}
}

func registerCoreMatchingRules(s *Schema) {
	rules := []*MatchingRule{
		{OID: MatchingRuleOIDBitStringMatch, Names: []string{"bitStringMatch"}, Syntax: SyntaxOIDBitString},
		{OID: MatchingRuleOIDBooleanMatch, Names: []string{"booleanMatch"}, Syntax: SyntaxOIDBoolean},
		{OID: MatchingRuleOIDCaseExactMatch, Names: []string{"caseExactMatch"}, Syntax: SyntaxOIDDirectoryString},
		{OID: MatchingRuleOIDCaseExactOrderingMatch, Names: []string{"caseExactOrderingMatch"}, Syntax: SyntaxOIDDirectoryString},
		{OID: MatchingRuleOIDCaseExactSubstringsMatch, Names: []string{"caseExactSubstringsMatch"}, Syntax: SyntaxOIDDirectoryString},
		{OID: MatchingRuleOIDCaseIgnoreMatch, Names: []string{"caseIgnoreMatch"}, Syntax: SyntaxOIDDirectoryString},
		{OID: MatchingRuleOIDCaseIgnoreOrderingMatch, Names: []string{"caseIgnoreOrderingMatch"}, Syntax: SyntaxOIDDirectoryString},
		{OID: MatchingRuleOIDCaseIgnoreSubstringsMatch, Names: []string{"caseIgnoreSubstringsMatch"}, Syntax: SyntaxOIDDirectoryString},
		{OID: MatchingRuleOIDDistinguishedNameMatch, Names: []string{"distinguishedNameMatch"}, Syntax: SyntaxOIDDN},
		{OID: MatchingRuleOIDGeneralizedTimeMatch, Names: []string{"generalizedTimeMatch"}, Syntax: SyntaxOIDGeneralizedTime},
		{OID: MatchingRuleOIDGeneralizedTimeOrderingMatch, Names: []string{"generalizedTimeOrderingMatch"}, Syntax: SyntaxOIDGeneralizedTime},
		{OID: MatchingRuleOIDIntegerMatch, Names: []string{"integerMatch"}, Syntax: SyntaxOIDInteger},
		{OID: MatchingRuleOIDIntegerOrderingMatch, Names: []string{"integerOrderingMatch"}, Syntax: SyntaxOIDInteger},
		{OID: MatchingRuleOIDNumericStringMatch, Names: []string{"numericStringMatch"}, Syntax: SyntaxOIDNumericString},
		{OID: MatchingRuleOIDNumericStringOrderingMatch, Names: []string{"numericStringOrderingMatch"}, Syntax: SyntaxOIDNumericString},
		{OID: MatchingRuleOIDNumericStringSubstringsMatch, Names: []string{"numericStringSubstringsMatch"}, Syntax: SyntaxOIDNumericString},
		{OID: MatchingRuleOIDObjectIdentifierMatch, Names: []string{"objectIdentifierMatch"}, Syntax: SyntaxOIDOID},
		{OID: MatchingRuleOIDOctetStringMatch, Names: []string{"octetStringMatch"}, Syntax: SyntaxOIDOctetString},
		{OID: MatchingRuleOIDOctetStringOrderingMatch, Names: []string{"octetStringOrderingMatch"}, Syntax: SyntaxOIDOctetString},
		{OID: MatchingRuleOIDTelephoneNumberMatch, Names: []string{"telephoneNumberMatch"}, Syntax: SyntaxOIDTelephoneNumber},
		{OID: MatchingRuleOIDTelephoneNumberSubstringsMatch, Names: []string{"telephoneNumberSubstringsMatch"}, Syntax: SyntaxOIDTelephoneNumber},
		{OID: MatchingRuleOIDUniqueIdentifierMatch, Names: []string{"uniqueIdentifierMatch"}, Syntax: SyntaxOIDBitString},
	}

	for _, rule := range rules {
		s.MatchingRules[rule.OID] = rule
		for _, name := range rule.Names {
			s.MatchingRules[name] = rule
		}
	}
}

func registerCoreAttributeTypes(s *Schema) {
	attributes := []*AttributeType{
		{
			OID: AttributeOIDBusinessCategory, Names: []string{"businessCategory"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDC, Names: []string{"c", "countryName"},
			SuperType: "name", Syntax: SyntaxOIDCountryString, SingleValue: true,
		},
		{
			OID: AttributeOIDCN, Names: []string{"cn", "commonName"},
			SuperType: "name", Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDDC, Names: []string{"dc", "domainComponent"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDIA5String, SingleValue: true,
		},
		{
			OID: AttributeOIDDescription, Names: []string{"description"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDDestinationIndicator, Names: []string{"destinationIndicator"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDPrintableString,
		},
		{
			OID: AttributeOIDDistinguishedName, Names: []string{"distinguishedName"},
			Equality: MatchingRuleOIDDistinguishedNameMatch, Syntax: SyntaxOIDDN,
		},
		{
			OID: AttributeOIDDNQualifier, Names: []string{"dnQualifier"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Ordering: MatchingRuleOIDCaseIgnoreOrderingMatch,
			Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch, Syntax: SyntaxOIDPrintableString,
		},
		{
			OID: AttributeOIDEnhancedSearchGuide, Names: []string{"enhancedSearchGuide"},
			Syntax: SyntaxOIDEnhancedGuide,
		},
		{
			OID: AttributeOIDFacsimileTelephoneNumber, Names: []string{"facsimileTelephoneNumber"},
			Syntax: SyntaxOIDFacsimileTelephoneNumber,
		},
		{
			OID: AttributeOIDGenerationQualifier, Names: []string{"generationQualifier"},
			SuperType: "name", Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDGivenName, Names: []string{"givenName"},
			SuperType: "name", Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDHouseIdentifier, Names: []string{"houseIdentifier"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDInitials, Names: []string{"initials"},
			SuperType: "name", Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDInternationalISDNNumber, Names: []string{"internationaliSDNNumber"},
			Equality: MatchingRuleOIDNumericStringMatch, Substring: MatchingRuleOIDNumericStringSubstringsMatch,
			Syntax: SyntaxOIDNumericString,
		},
		{
			OID: AttributeOIDL, Names: []string{"l", "localityName"},
			SuperType: "name", Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDMember, Names: []string{"member"},
			SuperType: "distinguishedName", Syntax: SyntaxOIDDN,
		},
		{
			OID: AttributeOIDName, Names: []string{"name"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDO, Names: []string{"o", "organizationName"},
			SuperType: "name", Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDOU, Names: []string{"ou", "organizationalUnitName"},
			SuperType: "name", Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDOwner, Names: []string{"owner"},
			SuperType: "distinguishedName", Syntax: SyntaxOIDDN,
		},
		{
			OID: AttributeOIDPhysicalDeliveryOfficeName, Names: []string{"physicalDeliveryOfficeName"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDPostalAddress, Names: []string{"postalAddress"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDPostalAddress,
		},
		{
			OID: AttributeOIDPostalCode, Names: []string{"postalCode"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDPostOfficeBox, Names: []string{"postOfficeBox"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDPreferredDeliveryMethod, Names: []string{"preferredDeliveryMethod"},
			Syntax: SyntaxOIDDirectoryString, SingleValue: true,
		},
		{
			OID: AttributeOIDRegisteredAddress, Names: []string{"registeredAddress"},
			SuperType: "postalAddress", Syntax: SyntaxOIDPostalAddress,
		},
		{
			OID: AttributeOIDRoleOccupant, Names: []string{"roleOccupant"},
			SuperType: "distinguishedName", Syntax: SyntaxOIDDN,
		},
		{
			OID: AttributeOIDSearchGuide, Names: []string{"searchGuide"},
			Syntax: SyntaxOIDGuide,
		},
		{
			OID: AttributeOIDSeeAlso, Names: []string{"seeAlso"},
			SuperType: "distinguishedName", Syntax: SyntaxOIDDN,
		},
		{
			OID: AttributeOIDSerialNumber, Names: []string{"serialNumber"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDPrintableString,
		},
		{
			OID: AttributeOIDSN, Names: []string{"sn", "surname"},
			SuperType: "name", Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDST, Names: []string{"st", "stateOrProvinceName"},
			SuperType: "name", Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDStreet, Names: []string{"street", "streetAddress"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDTelephoneNumber, Names: []string{"telephoneNumber"},
			Equality: MatchingRuleOIDTelephoneNumberMatch, Substring: MatchingRuleOIDTelephoneNumberSubstringsMatch,
			Syntax: SyntaxOIDTelephoneNumber,
		},
		{
			OID: AttributeOIDTeletexTerminalIdentifier, Names: []string{"teletexTerminalIdentifier"},
			Syntax: SyntaxOIDTeletexTerminalIdentifier,
		},
		{
			OID: AttributeOIDTelexNumber, Names: []string{"telexNumber"},
			Syntax: SyntaxOIDTelexNumber,
		},
		{
			OID: AttributeOIDTitle, Names: []string{"title"},
			SuperType: "name", Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDUID, Names: []string{"uid", "userid"},
			Equality: MatchingRuleOIDCaseIgnoreMatch, Substring: MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax: SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDUniqueMember, Names: []string{"uniqueMember"},
			Equality: MatchingRuleOIDDistinguishedNameMatch, Syntax: SyntaxOIDNameAndOptionalUID,
		},
		{
			OID: AttributeOIDUserPassword, Names: []string{"userPassword"},
			Equality: MatchingRuleOIDOctetStringMatch, Syntax: SyntaxOIDOctetString,
		},
		{
			OID: AttributeOIDX121Address, Names: []string{"x121Address"},
			Equality: MatchingRuleOIDNumericStringMatch, Substring: MatchingRuleOIDNumericStringSubstringsMatch,
			Syntax: SyntaxOIDNumericString,
		},
		{
			OID: AttributeOIDX500UniqueIdentifier, Names: []string{"x500UniqueIdentifier"},
			Equality: MatchingRuleOIDBitStringMatch, Syntax: SyntaxOIDBitString,
		},
	}

	for _, attr := range attributes {
		s.AttributeTypes[attr.OID] = attr
		for _, name := range attr.Names {
			s.AttributeTypes[name] = attr
		}
	}
}

func registerCoreObjectClasses(s *Schema) {
	classes := []*ObjectClass{
		{
			OID: ObjectClassOIDTop, Names: []string{"top"},
			Kind: ObjectClassKindAbstract, Must: []string{"objectClass"},
		},
		{
			OID: ObjectClassOIDAlias, Names: []string{"alias"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"aliasedObjectName"},
		},
		{
			OID: ObjectClassOIDCountry, Names: []string{"country"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"c"},
			May:  []string{"searchGuide", "description"},
		},
		{
			OID: ObjectClassOIDLocality, Names: []string{"locality"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"l"},
			May:  []string{"street", "seeAlso", "searchGuide", "st", "l", "description"},
		},
		{
			OID: ObjectClassOIDOrganization, Names: []string{"organization"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"o"},
			May: []string{"userPassword", "searchGuide", "seeAlso", "businessCategory",
				"x121Address", "registeredAddress", "destinationIndicator",
				"preferredDeliveryMethod", "telexNumber", "teletexTerminalIdentifier",
				"telephoneNumber", "internationaliSDNNumber", "facsimileTelephoneNumber",
				"street", "postOfficeBox", "postalCode", "postalAddress",
				"physicalDeliveryOfficeName", "st", "l", "description"},
		},
		{
			OID: ObjectClassOIDOrganizationalUnit, Names: []string{"organizationalUnit"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"ou"},
			May: []string{"businessCategory", "description", "destinationIndicator",
				"facsimileTelephoneNumber", "internationaliSDNNumber", "l",
				"physicalDeliveryOfficeName", "postalAddress", "postalCode",
				"postOfficeBox", "preferredDeliveryMethod", "registeredAddress",
				"searchGuide", "seeAlso", "st", "street", "telephoneNumber",
				"teletexTerminalIdentifier", "telexNumber", "userPassword",
				"x121Address"},
		},
		{
			OID: ObjectClassOIDPerson, Names: []string{"person"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"sn", "cn"},
			May:  []string{"userPassword", "telephoneNumber", "seeAlso", "description"},
		},
		{
			OID: ObjectClassOIDOrganizationalPerson, Names: []string{"organizationalPerson"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"person"},
			May: []string{"title", "x121Address", "registeredAddress", "destinationIndicator",
				"preferredDeliveryMethod", "telexNumber", "teletexTerminalIdentifier",
				"telephoneNumber", "internationaliSDNNumber", "facsimileTelephoneNumber",
				"street", "postOfficeBox", "postalCode", "postalAddress",
				"physicalDeliveryOfficeName", "ou", "st", "l"},
		},
		{
			OID: ObjectClassOIDOrganizationalRole, Names: []string{"organizationalRole"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"cn"},
			May: []string{"x121Address", "registeredAddress", "destinationIndicator",
				"preferredDeliveryMethod", "telexNumber", "teletexTerminalIdentifier",
				"telephoneNumber", "internationaliSDNNumber", "facsimileTelephoneNumber",
				"seeAlso", "roleOccupant", "preferredDeliveryMethod", "street",
				"postOfficeBox", "postalCode", "postalAddress",
				"physicalDeliveryOfficeName", "ou", "st", "l", "description"},
		},
		{
			OID: ObjectClassOIDGroupOfNames, Names: []string{"groupOfNames"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"member", "cn"},
			May:  []string{"businessCategory", "seeAlso", "owner", "ou", "o", "description"},
		},
		{
			OID: ObjectClassOIDResidentialPerson, Names: []string{"residentialPerson"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"person"},
			Must: []string{"l"},
			May: []string{"businessCategory", "x121Address", "registeredAddress",
				"destinationIndicator", "preferredDeliveryMethod", "telexNumber",
				"teletexTerminalIdentifier", "telephoneNumber", "internationaliSDNNumber",
				"facsimileTelephoneNumber", "preferredDeliveryMethod", "street",
				"postOfficeBox", "postalCode", "postalAddress",
				"physicalDeliveryOfficeName", "st", "l"},
		},
		{
			OID: ObjectClassOIDApplicationProcess, Names: []string{"applicationProcess"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"cn"},
			May:  []string{"seeAlso", "ou", "l", "description"},
		},
		{
			OID: ObjectClassOIDDevice, Names: []string{"device"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"cn"},
			May:  []string{"serialNumber", "seeAlso", "owner", "ou", "o", "l", "description"},
		},
		{
			OID: ObjectClassOIDGroupOfUniqueNames, Names: []string{"groupOfUniqueNames"},
			Kind: ObjectClassKindStructural, SuperClasses: []string{"top"},
			Must: []string{"uniqueMember", "cn"},
			May:  []string{"businessCategory", "seeAlso", "owner", "ou", "o", "description"},
		},
		{
			OID: ObjectClassOIDDCObject, Names: []string{"dcObject"},
			Kind: ObjectClassKindAuxiliary, SuperClasses: []string{"top"},
			Must: []string{"dc"},
		},
		{
			OID: ObjectClassOIDUIDObject, Names: []string{"uidObject"},
			Kind: ObjectClassKindAuxiliary, SuperClasses: []string{"top"},
			Must: []string{"uid"},
		},
	}

	for _, class := range classes {
		s.ObjectClasses[class.OID] = class
		for _, name := range class.Names {
			s.ObjectClasses[name] = class
		}
	}
}
