package mod_ldap

// Operational Attribute OIDs (RFC 4512, RFC 4530)
// These attributes are maintained automatically by the LDAP server.
const (
	// AttributeOIDCreateTimestamp is the standard operational attribute for entry creation time (RFC 4512).
	AttributeOIDCreateTimestamp = "2.5.18.1"
	// AttributeOIDModifyTimestamp is the standard operational attribute for entry modification time (RFC 4512).
	AttributeOIDModifyTimestamp = "2.5.18.2"
	// AttributeOIDCreatorsName is the standard operational attribute for creator DN (RFC 4512).
	AttributeOIDCreatorsName = "2.5.18.3"
	// AttributeOIDModifiersName is the standard operational attribute for modifier DN (RFC 4512).
	AttributeOIDModifiersName = "2.5.18.4"
	// AttributeOIDSubschemaSubentry is the standard operational attribute for schema reference (RFC 4512).
	AttributeOIDSubschemaSubentry = "2.5.18.10"
	// AttributeOIDStructuralObjectClass is the standard operational attribute for structural object class (RFC 4512).
	AttributeOIDStructuralObjectClass = "2.5.21.9"

	// AttributeOIDEntryUUID is the operational attribute for entry UUID (RFC 4530).
	AttributeOIDEntryUUID = "1.3.6.1.1.16.4"

	// AttributeOIDEntryDN is the operational attribute for entry DN (implementation-specific).
	AttributeOIDEntryDN = "1.3.6.1.1.20"
	// AttributeOIDHasSubordinates is the operational attribute indicating if entry has children.
	AttributeOIDHasSubordinates = "2.5.18.9"
	// AttributeOIDSubordinateCount is the operational attribute for number of children (Novell).
	AttributeOIDSubordinateCount = "2.16.840.1.113719.1.1.4.1.6"
	// AttributeOIDEntryCSN is the operational attribute for change sequence number (OpenLDAP).
	AttributeOIDEntryCSN = "1.3.6.1.4.1.4203.666.1.7"
	// AttributeOIDContextCSN is the operational attribute for context CSN (OpenLDAP).
	AttributeOIDContextCSN = "1.3.6.1.4.1.4203.666.1.25"
)

// Attribute Usage Types for operational attributes.
const (
	AttributeUsageDirectoryOperationLong   = "directoryOperation"
	AttributeUsageDistributedOperationLong = "distributedOperation"
	AttributeUsageDSAOperationLong         = "dSAOperation"
)
