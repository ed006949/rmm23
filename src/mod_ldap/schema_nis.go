package mod_ldap

// NISSchema returns the RFC 2307 Network Information Service LDAP schema.
func NISSchema() *Schema {
	schema := NewSchema()

	// Register NIS attribute types
	registerNISAttributeTypes(schema)

	// Register NIS object classes
	registerNISObjectClasses(schema)

	return schema
}

func registerNISAttributeTypes(s *Schema) {
	attributes := []*AttributeType{
		{
			OID: AttributeOIDUIDNumber, Names: []string{"uidNumber"},
			Description: "An integer uniquely identifying a user in an administrative domain",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDGIDNumber, Names: []string{"gidNumber"},
			Description: "An integer uniquely identifying a group in an administrative domain",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDGecos, Names: []string{"gecos"},
			Description: "The GECOS field; the common name",
			Equality:    MatchingRuleOIDCaseIgnoreIA5Match,
			Substring:   MatchingRuleOIDCaseIgnoreIA5SubstringsMatch,
			Syntax:      SyntaxOIDIA5String,
			SingleValue: true,
		},
		{
			OID: AttributeOIDHomeDirectory, Names: []string{"homeDirectory"},
			Description: "The absolute path to the home directory",
			Equality:    MatchingRuleOIDCaseExactMatch,
			Syntax:      SyntaxOIDIA5String,
			SingleValue: true,
		},
		{
			OID: AttributeOIDLoginShell, Names: []string{"loginShell"},
			Description: "The path to the login shell",
			Equality:    MatchingRuleOIDCaseExactMatch,
			Syntax:      SyntaxOIDIA5String,
			SingleValue: true,
		},
		{
			OID: AttributeOIDShadowLastChange, Names: []string{"shadowLastChange"},
			Description: "Days since Jan 1, 1970 that password was last changed",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDShadowMin, Names: []string{"shadowMin"},
			Description: "Days before password may be changed",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDShadowMax, Names: []string{"shadowMax"},
			Description: "Days after which password must be changed",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDShadowWarning, Names: []string{"shadowWarning"},
			Description: "Days before password is to expire that user is warned",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDShadowInactive, Names: []string{"shadowInactive"},
			Description: "Days after password expires that account is disabled",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDShadowExpire, Names: []string{"shadowExpire"},
			Description: "Days since Jan 1, 1970 that account is disabled",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDShadowFlag, Names: []string{"shadowFlag"},
			Description: "Reserved field",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDMemberUID, Names: []string{"memberUid"},
			Description: "Member user identifier of a group",
			Equality:    MatchingRuleOIDCaseExactMatch,
			Substring:   MatchingRuleOIDCaseExactSubstringsMatch,
			Syntax:      SyntaxOIDIA5String,
		},
		{
			OID: AttributeOIDMemberNisNetgroup, Names: []string{"memberNisNetgroup"},
			Description: "Member NIS netgroup",
			Equality:    MatchingRuleOIDCaseExactMatch,
			Substring:   MatchingRuleOIDCaseExactSubstringsMatch,
			Syntax:      SyntaxOIDIA5String,
		},
		{
			OID: AttributeOIDNisNetgroupTriple, Names: []string{"nisNetgroupTriple"},
			Description: "Netgroup triple (host,user,domain)",
			Syntax:      SyntaxOIDIA5String,
		},
		{
			OID: AttributeOIDIpServicePort, Names: []string{"ipServicePort"},
			Description: "Port number for IP service",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDIpServiceProtocol, Names: []string{"ipServiceProtocol"},
			Description: "IP protocol number or name",
			Equality:    MatchingRuleOIDCaseIgnoreIA5Match,
			Syntax:      SyntaxOIDIA5String,
		},
		{
			OID: AttributeOIDIpProtocolNumber, Names: []string{"ipProtocolNumber"},
			Description: "IP protocol number",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDOncRpcNumber, Names: []string{"oncRpcNumber"},
			Description: "ONC RPC number",
			Equality:    MatchingRuleOIDIntegerMatch,
			Syntax:      SyntaxOIDInteger,
			SingleValue: true,
		},
		{
			OID: AttributeOIDIpHostNumber, Names: []string{"ipHostNumber"},
			Description: "IP address as a dotted decimal",
			Equality:    MatchingRuleOIDCaseIgnoreIA5Match,
			Syntax:      SyntaxOIDIA5String,
		},
		{
			OID: AttributeOIDIpNetworkNumber, Names: []string{"ipNetworkNumber"},
			Description: "IP network as a dotted decimal",
			Equality:    MatchingRuleOIDCaseIgnoreIA5Match,
			Substring:   MatchingRuleOIDCaseIgnoreIA5SubstringsMatch,
			Syntax:      SyntaxOIDIA5String,
			SingleValue: true,
		},
		{
			OID: AttributeOIDIpNetmaskNumber, Names: []string{"ipNetmaskNumber"},
			Description: "IP netmask as a dotted decimal",
			Equality:    MatchingRuleOIDCaseIgnoreIA5Match,
			Substring:   MatchingRuleOIDCaseIgnoreIA5SubstringsMatch,
			Syntax:      SyntaxOIDIA5String,
			SingleValue: true,
		},
		{
			OID: AttributeOIDMacAddress, Names: []string{"macAddress"},
			Description: "MAC address in hexadecimal format",
			Equality:    MatchingRuleOIDCaseIgnoreIA5Match,
			Syntax:      SyntaxOIDIA5String,
		},
		{
			OID: AttributeOIDBootParameter, Names: []string{"bootParameter"},
			Description: "Boot parameter for network boot",
			Equality:    MatchingRuleOIDCaseExactMatch,
			Syntax:      SyntaxOIDIA5String,
		},
		{
			OID: AttributeOIDBootFile, Names: []string{"bootFile"},
			Description: "Boot image filename",
			Equality:    MatchingRuleOIDCaseExactMatch,
			Substring:   MatchingRuleOIDCaseExactSubstringsMatch,
			Syntax:      SyntaxOIDIA5String,
		},
		{
			OID: AttributeOIDNisMapName, Names: []string{"nisMapName"},
			Description: "Name of a generic NIS map",
			Equality:    MatchingRuleOIDCaseIgnoreMatch,
			Syntax:      SyntaxOIDDirectoryString,
		},
		{
			OID: AttributeOIDNisMapEntry, Names: []string{"nisMapEntry"},
			Description: "Single entry of a NIS map",
			Equality:    MatchingRuleOIDCaseExactMatch,
			Substring:   MatchingRuleOIDCaseExactSubstringsMatch,
			Syntax:      SyntaxOIDIA5String,
			SingleValue: true,
		},
	}

	for _, attr := range attributes {
		s.AttributeTypes[attr.OID] = attr
		for _, name := range attr.Names {
			s.AttributeTypes[name] = attr
		}
	}
}

func registerNISObjectClasses(s *Schema) {
	classes := []*ObjectClass{
		{
			OID:          ObjectClassOIDPosixAccount,
			Names:        []string{"posixAccount"},
			Description:  "Abstraction of an account with POSIX attributes",
			Kind:         ObjectClassKindAuxiliary,
			SuperClasses: []string{"top"},
			Must:         []string{"cn", "uid", "uidNumber", "gidNumber", "homeDirectory"},
			May:          []string{"userPassword", "loginShell", "gecos", "description"},
		},
		{
			OID:          ObjectClassOIDShadowAccount,
			Names:        []string{"shadowAccount"},
			Description:  "Additional attributes for shadow password",
			Kind:         ObjectClassKindAuxiliary,
			SuperClasses: []string{"top"},
			Must:         []string{"uid"},
			May: []string{"userPassword", "shadowLastChange", "shadowMin", "shadowMax",
				"shadowWarning", "shadowInactive", "shadowExpire", "shadowFlag", "description"},
		},
		{
			OID:          ObjectClassOIDPosixGroup,
			Names:        []string{"posixGroup"},
			Description:  "Abstraction of a group of accounts",
			Kind:         ObjectClassKindStructural,
			SuperClasses: []string{"top"},
			Must:         []string{"cn", "gidNumber"},
			May:          []string{"userPassword", "memberUid", "description"},
		},
		{
			OID:          ObjectClassOIDIpService,
			Names:        []string{"ipService"},
			Description:  "Abstraction of an Internet Protocol service",
			Kind:         ObjectClassKindStructural,
			SuperClasses: []string{"top"},
			Must:         []string{"cn", "ipServicePort", "ipServiceProtocol"},
			May:          []string{"description"},
		},
		{
			OID:          ObjectClassOIDIpProtocol,
			Names:        []string{"ipProtocol"},
			Description:  "Abstraction of an IP protocol",
			Kind:         ObjectClassKindStructural,
			SuperClasses: []string{"top"},
			Must:         []string{"cn", "ipProtocolNumber"},
			May:          []string{"description"},
		},
		{
			OID:          ObjectClassOIDOncRpc,
			Names:        []string{"oncRpc"},
			Description:  "Abstraction of an Open Network Computing (ONC) Remote Procedure Call (RPC) binding",
			Kind:         ObjectClassKindStructural,
			SuperClasses: []string{"top"},
			Must:         []string{"cn", "oncRpcNumber"},
			May:          []string{"description"},
		},
		{
			OID:          ObjectClassOIDIpHost,
			Names:        []string{"ipHost"},
			Description:  "Abstraction of a host, an IP device",
			Kind:         ObjectClassKindAuxiliary,
			SuperClasses: []string{"top"},
			Must:         []string{"cn", "ipHostNumber"},
			May:          []string{"l", "description", "manager"},
		},
		{
			OID:          ObjectClassOIDIpNetwork,
			Names:        []string{"ipNetwork"},
			Description:  "Abstraction of a network",
			Kind:         ObjectClassKindStructural,
			SuperClasses: []string{"top"},
			Must:         []string{"cn", "ipNetworkNumber"},
			May:          []string{"ipNetmaskNumber", "l", "description", "manager"},
		},
		{
			OID:          ObjectClassOIDNisNetgroup,
			Names:        []string{"nisNetgroup"},
			Description:  "Abstraction of a netgroup",
			Kind:         ObjectClassKindStructural,
			SuperClasses: []string{"top"},
			Must:         []string{"cn"},
			May:          []string{"nisNetgroupTriple", "memberNisNetgroup", "description"},
		},
		{
			OID:          ObjectClassOIDNisMap,
			Names:        []string{"nisMap"},
			Description:  "A generic abstraction of a NIS map",
			Kind:         ObjectClassKindStructural,
			SuperClasses: []string{"top"},
			Must:         []string{"nisMapName"},
			May:          []string{"description"},
		},
		{
			OID:          ObjectClassOIDNisObject,
			Names:        []string{"nisObject"},
			Description:  "An entry in a NIS map",
			Kind:         ObjectClassKindStructural,
			SuperClasses: []string{"top"},
			Must:         []string{"cn", "nisMapEntry", "nisMapName"},
			May:          []string{"description"},
		},
		{
			OID:          ObjectClassOIDIeee802Device,
			Names:        []string{"ieee802Device"},
			Description:  "A device with a MAC address",
			Kind:         ObjectClassKindAuxiliary,
			SuperClasses: []string{"top"},
			Must:         []string{"macAddress"},
		},
		{
			OID:          ObjectClassOIDBootableDevice,
			Names:        []string{"bootableDevice"},
			Description:  "A device that can be booted over the network",
			Kind:         ObjectClassKindAuxiliary,
			SuperClasses: []string{"top"},
			May:          []string{"bootFile", "bootParameter"},
		},
	}

	for _, class := range classes {
		s.ObjectClasses[class.OID] = class
		for _, name := range class.Names {
			s.ObjectClasses[name] = class
		}
	}
}
