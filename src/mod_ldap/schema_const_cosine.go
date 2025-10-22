package mod_ldap

// COSINE Schema Attribute OIDs (RFC 4524).
const (
	AttributeOIDAssociatedDomain     = "0.9.2342.19200300.100.1.37"
	AttributeOIDAssociatedName       = "0.9.2342.19200300.100.1.38"
	AttributeOIDBuildingName         = "0.9.2342.19200300.100.1.48"
	AttributeOIDCo                   = "0.9.2342.19200300.100.1.43"
	AttributeOIDDocumentAuthor       = "0.9.2342.19200300.100.1.14"
	AttributeOIDDocumentIdentifier   = "0.9.2342.19200300.100.1.11"
	AttributeOIDDocumentLocation     = "0.9.2342.19200300.100.1.15"
	AttributeOIDDocumentPublisher    = "0.9.2342.19200300.100.1.56"
	AttributeOIDDocumentTitle        = "0.9.2342.19200300.100.1.12"
	AttributeOIDDocumentVersion      = "0.9.2342.19200300.100.1.13"
	AttributeOIDDrink                = "0.9.2342.19200300.100.1.5"
	AttributeOIDHomePhone            = "0.9.2342.19200300.100.1.20"
	AttributeOIDHomePostalAddress    = "0.9.2342.19200300.100.1.39"
	AttributeOIDHost                 = "0.9.2342.19200300.100.1.9"
	AttributeOIDInfo                 = "0.9.2342.19200300.100.1.4"
	AttributeOIDMail                 = "0.9.2342.19200300.100.1.3"
	AttributeOIDManager              = "0.9.2342.19200300.100.1.10"
	AttributeOIDMobile               = "0.9.2342.19200300.100.1.41"
	AttributeOIDOrganizationalStatus = "0.9.2342.19200300.100.1.45"
	AttributeOIDPager                = "0.9.2342.19200300.100.1.42"
	AttributeOIDPersonalTitle        = "0.9.2342.19200300.100.1.40"
	AttributeOIDRoomNumber           = "0.9.2342.19200300.100.1.6"
	AttributeOIDSecretary            = "0.9.2342.19200300.100.1.21"
	AttributeOIDUniqueIdentifier     = "0.9.2342.19200300.100.1.44"
	AttributeOIDUserClass            = "0.9.2342.19200300.100.1.8"
)

// COSINE Schema Object Class OIDs (RFC 4524).
const (
	ObjectClassOIDAccount              = "0.9.2342.19200300.100.4.5"
	ObjectClassOIDDocument             = "0.9.2342.19200300.100.4.6"
	ObjectClassOIDDocumentSeries       = "0.9.2342.19200300.100.4.9"
	ObjectClassOIDDomain               = "0.9.2342.19200300.100.4.13"
	ObjectClassOIDDomainRelatedObject  = "0.9.2342.19200300.100.4.17"
	ObjectClassOIDFriendlyCountry      = "0.9.2342.19200300.100.4.18"
	ObjectClassOIDRFC822LocalPart      = "0.9.2342.19200300.100.4.14"
	ObjectClassOIDRoom                 = "0.9.2342.19200300.100.4.7"
	ObjectClassOIDSimpleSecurityObject = "0.9.2342.19200300.100.4.19"
)

// Additional matching rule OIDs used by COSINE schema.
const (
	MatchingRuleOIDCaseIgnoreIA5Match           = "1.3.6.1.4.1.1466.109.114.1"
	MatchingRuleOIDCaseIgnoreIA5SubstringsMatch = "1.3.6.1.4.1.1466.109.114.2"
)
