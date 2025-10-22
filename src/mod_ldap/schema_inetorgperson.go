package mod_ldap

// InetOrgPersonSchema returns the RFC 2798 inetOrgPerson LDAP schema.
func InetOrgPersonSchema() *Schema {
	schema := NewSchema()

	// Register inetOrgPerson attribute types
	registerInetOrgPersonAttributeTypes(schema)

	// Register inetOrgPerson object class
	registerInetOrgPersonObjectClasses(schema)

	return schema
}

func registerInetOrgPersonAttributeTypes(s *Schema) {
	attributes := []*AttributeType{
		{
			OID: AttributeOIDCarLicense, Names: []string{"carLicense"},
			Description: "vehicle license or registration plate",
			Equality:    MatchingRuleOIDCaseIgnoreMatch,
			Substring:   MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax:      SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDDepartmentNumber, Names: []string{"departmentNumber"},
			Description: "identifies a department within an organization",
			Equality:    MatchingRuleOIDCaseIgnoreMatch,
			Substring:   MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax:      SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDDisplayName, Names: []string{"displayName"},
			Description: "preferred name of a person to be used when displaying entries",
			Equality:    MatchingRuleOIDCaseIgnoreMatch,
			Substring:   MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax:      SyntaxOIDDirectoryString,
			SingleValue: true,
		},
		{
			OID: AttributeOIDEmployeeNumber, Names: []string{"employeeNumber"},
			Description: "numerically identifies an employee within an organization",
			Equality:    MatchingRuleOIDCaseIgnoreMatch,
			Substring:   MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax:      SyntaxOIDDirectoryString,
			SingleValue: true,
		},
		{
			OID: AttributeOIDEmployeeType, Names: []string{"employeeType"},
			Description: "type of employment for a person",
			Equality:    MatchingRuleOIDCaseIgnoreMatch,
			Substring:   MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax:      SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDJPEGPhoto, Names: []string{"jpegPhoto"},
			Description: "a JPEG image",
			Syntax:      SyntaxOIDJPEG,
		},
		{
			OID: AttributeOIDPreferredLanguage, Names: []string{"preferredLanguage"},
			Description: "preferred written or spoken language for a person",
			Equality:    MatchingRuleOIDCaseIgnoreMatch,
			Substring:   MatchingRuleOIDCaseIgnoreSubstringsMatch,
			Syntax:      SyntaxOIDDirectoryString,
			SingleValue: true,
		},
		{
			OID: AttributeOIDUserSMIMECertificate, Names: []string{"userSMIMECertificate"},
			Description: "PKCS#7 SignedData used to support S/MIME",
			Syntax:      SyntaxOIDOctetString,
		},
		{
			OID: AttributeOIDUserPKCS12, Names: []string{"userPKCS12"},
			Description: "PKCS #12 PFX PDU for exchange of personal identity information",
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

func registerInetOrgPersonObjectClasses(s *Schema) {
	classes := []*ObjectClass{
		{
			OID:          ObjectClassOIDInetOrgPerson,
			Names:        []string{"inetOrgPerson"},
			Description:  "RFC2798: Internet Organizational Person",
			Kind:         ObjectClassKindStructural,
			SuperClasses: []string{"organizationalPerson"},
			May: []string{
				"audio", "businessCategory", "carLicense", "departmentNumber",
				"displayName", "employeeNumber", "employeeType", "givenName",
				"homePhone", "homePostalAddress", "initials", "jpegPhoto",
				"labeledURI", "mail", "manager", "mobile", "o", "pager",
				"photo", "roomNumber", "secretary", "uid", "userCertificate",
				"x500uniqueIdentifier", "preferredLanguage", "userSMIMECertificate",
				"userPKCS12",
			},
		},
	}

	for _, class := range classes {
		s.ObjectClasses[class.OID] = class
		for _, name := range class.Names {
			s.ObjectClasses[name] = class
		}
	}
}
