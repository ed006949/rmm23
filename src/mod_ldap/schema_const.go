package mod_ldap

// Object Class Kinds.
const (
	ObjectClassKindStructural = "STRUCTURAL"
	ObjectClassKindAuxiliary  = "AUXILIARY"
	ObjectClassKindAbstract   = "ABSTRACT"
)

// Attribute Usage Types.
const (
	AttributeUsageUserApplications     = "userApplications"
	AttributeUsageDirectoryOperation   = "directoryOperation"
	AttributeUsageDistributedOperation = "distributedOperation"
	AttributeUsageDSAOperation         = "dSAOperation"
)

// Common attribute constraints.
const (
	MaxLengthDirectoryString256 = 256
)

// LDAP Syntax OIDs (RFC 4517).
const (
	SyntaxOIDAttributeTypeDescription    = "1.3.6.1.4.1.1466.115.121.1.3"
	SyntaxOIDBitString                   = "1.3.6.1.4.1.1466.115.121.1.6"
	SyntaxOIDBoolean                     = "1.3.6.1.4.1.1466.115.121.1.7"
	SyntaxOIDCountryString               = "1.3.6.1.4.1.1466.115.121.1.11"
	SyntaxOIDDN                          = "1.3.6.1.4.1.1466.115.121.1.12"
	SyntaxOIDDirectoryString             = "1.3.6.1.4.1.1466.115.121.1.15"
	SyntaxOIDDITContentRuleDescription   = "1.3.6.1.4.1.1466.115.121.1.16"
	SyntaxOIDDITStructureRuleDescription = "1.3.6.1.4.1.1466.115.121.1.17"
	SyntaxOIDEnhancedGuide               = "1.3.6.1.4.1.1466.115.121.1.21"
	SyntaxOIDFacsimileTelephoneNumber    = "1.3.6.1.4.1.1466.115.121.1.22"
	SyntaxOIDFax                         = "1.3.6.1.4.1.1466.115.121.1.23"
	SyntaxOIDGeneralizedTime             = "1.3.6.1.4.1.1466.115.121.1.24"
	SyntaxOIDGuide                       = "1.3.6.1.4.1.1466.115.121.1.25"
	SyntaxOIDIA5String                   = "1.3.6.1.4.1.1466.115.121.1.26"
	SyntaxOIDInteger                     = "1.3.6.1.4.1.1466.115.121.1.27"
	SyntaxOIDJPEG                        = "1.3.6.1.4.1.1466.115.121.1.28"
	SyntaxOIDMatchingRuleDescription     = "1.3.6.1.4.1.1466.115.121.1.30"
	SyntaxOIDMatchingRuleUseDescription  = "1.3.6.1.4.1.1466.115.121.1.31"
	SyntaxOIDNameAndOptionalUID          = "1.3.6.1.4.1.1466.115.121.1.34"
	SyntaxOIDNameFormDescription         = "1.3.6.1.4.1.1466.115.121.1.35"
	SyntaxOIDNumericString               = "1.3.6.1.4.1.1466.115.121.1.36"
	SyntaxOIDObjectClassDescription      = "1.3.6.1.4.1.1466.115.121.1.37"
	SyntaxOIDOctetString                 = "1.3.6.1.4.1.1466.115.121.1.40"
	SyntaxOIDOID                         = "1.3.6.1.4.1.1466.115.121.1.38"
	SyntaxOIDPostalAddress               = "1.3.6.1.4.1.1466.115.121.1.41"
	SyntaxOIDPrintableString             = "1.3.6.1.4.1.1466.115.121.1.44"
	SyntaxOIDTelephoneNumber             = "1.3.6.1.4.1.1466.115.121.1.50"
	SyntaxOIDTeletexTerminalIdentifier   = "1.3.6.1.4.1.1466.115.121.1.51"
	SyntaxOIDTelexNumber                 = "1.3.6.1.4.1.1466.115.121.1.52"
)

// Matching Rule OIDs (RFC 4517).
const (
	MatchingRuleOIDBitStringMatch                 = "2.5.13.16"
	MatchingRuleOIDBooleanMatch                   = "2.5.13.13"
	MatchingRuleOIDCaseExactMatch                 = "2.5.13.5"
	MatchingRuleOIDCaseExactOrderingMatch         = "2.5.13.6"
	MatchingRuleOIDCaseExactSubstringsMatch       = "2.5.13.7"
	MatchingRuleOIDCaseIgnoreMatch                = "2.5.13.2"
	MatchingRuleOIDCaseIgnoreOrderingMatch        = "2.5.13.3"
	MatchingRuleOIDCaseIgnoreSubstringsMatch      = "2.5.13.4"
	MatchingRuleOIDDistinguishedNameMatch         = "2.5.13.1"
	MatchingRuleOIDGeneralizedTimeMatch           = "2.5.13.27"
	MatchingRuleOIDGeneralizedTimeOrderingMatch   = "2.5.13.28"
	MatchingRuleOIDIntegerMatch                   = "2.5.13.14"
	MatchingRuleOIDIntegerOrderingMatch           = "2.5.13.15"
	MatchingRuleOIDNumericStringMatch             = "2.5.13.8"
	MatchingRuleOIDNumericStringOrderingMatch     = "2.5.13.9"
	MatchingRuleOIDNumericStringSubstringsMatch   = "2.5.13.10"
	MatchingRuleOIDObjectIdentifierMatch          = "2.5.13.0"
	MatchingRuleOIDOctetStringMatch               = "2.5.13.17"
	MatchingRuleOIDOctetStringOrderingMatch       = "2.5.13.18"
	MatchingRuleOIDTelephoneNumberMatch           = "2.5.13.20"
	MatchingRuleOIDTelephoneNumberSubstringsMatch = "2.5.13.21"
	MatchingRuleOIDUniqueIdentifierMatch          = "2.5.13.23"
)

// Core Schema Attribute OIDs (RFC 4519).
const (
	AttributeOIDBusinessCategory           = "2.5.4.15"
	AttributeOIDC                          = "2.5.4.6"                    // countryName
	AttributeOIDCN                         = "2.5.4.3"                    // commonName
	AttributeOIDDC                         = "0.9.2342.19200300.100.1.25" // domainComponent
	AttributeOIDDescription                = "2.5.4.13"
	AttributeOIDDestinationIndicator       = "2.5.4.27"
	AttributeOIDDistinguishedName          = "2.5.4.49"
	AttributeOIDDNQualifier                = "2.5.4.46"
	AttributeOIDEnhancedSearchGuide        = "2.5.4.47"
	AttributeOIDFacsimileTelephoneNumber   = "2.5.4.23"
	AttributeOIDGenerationQualifier        = "2.5.4.44"
	AttributeOIDGivenName                  = "2.5.4.42"
	AttributeOIDHouseIdentifier            = "2.5.4.51"
	AttributeOIDInitials                   = "2.5.4.43"
	AttributeOIDInternationalISDNNumber    = "2.5.4.25"
	AttributeOIDL                          = "2.5.4.7" // localityName
	AttributeOIDMember                     = "2.5.4.31"
	AttributeOIDName                       = "2.5.4.41"
	AttributeOIDO                          = "2.5.4.10" // organizationName
	AttributeOIDOU                         = "2.5.4.11" // organizationalUnitName
	AttributeOIDOwner                      = "2.5.4.32"
	AttributeOIDPhysicalDeliveryOfficeName = "2.5.4.19"
	AttributeOIDPostalAddress              = "2.5.4.16"
	AttributeOIDPostalCode                 = "2.5.4.17"
	AttributeOIDPostOfficeBox              = "2.5.4.18"
	AttributeOIDPreferredDeliveryMethod    = "2.5.4.28"
	AttributeOIDRegisteredAddress          = "2.5.4.26"
	AttributeOIDRoleOccupant               = "2.5.4.33"
	AttributeOIDSearchGuide                = "2.5.4.14"
	AttributeOIDSeeAlso                    = "2.5.4.34"
	AttributeOIDSerialNumber               = "2.5.4.5"
	AttributeOIDSN                         = "2.5.4.4" // surname
	AttributeOIDST                         = "2.5.4.8" // stateOrProvinceName
	AttributeOIDStreet                     = "2.5.4.9" // streetAddress
	AttributeOIDTelephoneNumber            = "2.5.4.20"
	AttributeOIDTeletexTerminalIdentifier  = "2.5.4.22"
	AttributeOIDTelexNumber                = "2.5.4.21"
	AttributeOIDTitle                      = "2.5.4.12"
	AttributeOIDUID                        = "0.9.2342.19200300.100.1.1" // userid
	AttributeOIDUniqueMember               = "2.5.4.50"
	AttributeOIDUserPassword               = "2.5.4.35"
	AttributeOIDX121Address                = "2.5.4.24"
	AttributeOIDX500UniqueIdentifier       = "2.5.4.45"
)

// Core Schema Object Class OIDs (RFC 4519).
const (
	ObjectClassOIDTop                  = "2.5.6.0"
	ObjectClassOIDAlias                = "2.5.6.1"
	ObjectClassOIDCountry              = "2.5.6.2"
	ObjectClassOIDLocality             = "2.5.6.3"
	ObjectClassOIDOrganization         = "2.5.6.4"
	ObjectClassOIDOrganizationalUnit   = "2.5.6.5"
	ObjectClassOIDPerson               = "2.5.6.6"
	ObjectClassOIDOrganizationalPerson = "2.5.6.7"
	ObjectClassOIDOrganizationalRole   = "2.5.6.8"
	ObjectClassOIDGroupOfNames         = "2.5.6.9"
	ObjectClassOIDResidentialPerson    = "2.5.6.10"
	ObjectClassOIDApplicationProcess   = "2.5.6.11"
	ObjectClassOIDDevice               = "2.5.6.14"
	ObjectClassOIDGroupOfUniqueNames   = "2.5.6.17"
	ObjectClassOIDDCObject             = "1.3.6.1.4.1.1466.344"
	ObjectClassOIDUIDObject            = "1.3.6.1.1.3.1"
)
