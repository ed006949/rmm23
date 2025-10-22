package mod_ldap

// NIS Schema Base OID (RFC 2307).
const (
	NISSchemaBase = "1.3.6.1.1.1" // nisSchema
)

// NIS Schema Attribute OIDs (RFC 2307).
const (
	AttributeOIDUIDNumber         = "1.3.6.1.1.1.1.0"  // nisSchema.1.0
	AttributeOIDGIDNumber         = "1.3.6.1.1.1.1.1"  // nisSchema.1.1
	AttributeOIDGecos             = "1.3.6.1.1.1.1.2"  // nisSchema.1.2
	AttributeOIDHomeDirectory     = "1.3.6.1.1.1.1.3"  // nisSchema.1.3
	AttributeOIDLoginShell        = "1.3.6.1.1.1.1.4"  // nisSchema.1.4
	AttributeOIDShadowLastChange  = "1.3.6.1.1.1.1.5"  // nisSchema.1.5
	AttributeOIDShadowMin         = "1.3.6.1.1.1.1.6"  // nisSchema.1.6
	AttributeOIDShadowMax         = "1.3.6.1.1.1.1.7"  // nisSchema.1.7
	AttributeOIDShadowWarning     = "1.3.6.1.1.1.1.8"  // nisSchema.1.8
	AttributeOIDShadowInactive    = "1.3.6.1.1.1.1.9"  // nisSchema.1.9
	AttributeOIDShadowExpire      = "1.3.6.1.1.1.1.10" // nisSchema.1.10
	AttributeOIDShadowFlag        = "1.3.6.1.1.1.1.11" // nisSchema.1.11
	AttributeOIDMemberUID         = "1.3.6.1.1.1.1.12" // nisSchema.1.12
	AttributeOIDMemberNisNetgroup = "1.3.6.1.1.1.1.13" // nisSchema.1.13
	AttributeOIDNisNetgroupTriple = "1.3.6.1.1.1.1.14" // nisSchema.1.14
	AttributeOIDIpServicePort     = "1.3.6.1.1.1.1.15" // nisSchema.1.15
	AttributeOIDIpServiceProtocol = "1.3.6.1.1.1.1.16" // nisSchema.1.16
	AttributeOIDIpProtocolNumber  = "1.3.6.1.1.1.1.17" // nisSchema.1.17
	AttributeOIDOncRpcNumber      = "1.3.6.1.1.1.1.18" // nisSchema.1.18
	AttributeOIDIpHostNumber      = "1.3.6.1.1.1.1.19" // nisSchema.1.19
	AttributeOIDIpNetworkNumber   = "1.3.6.1.1.1.1.20" // nisSchema.1.20
	AttributeOIDIpNetmaskNumber   = "1.3.6.1.1.1.1.21" // nisSchema.1.21
	AttributeOIDMacAddress        = "1.3.6.1.1.1.1.22" // nisSchema.1.22
	AttributeOIDBootParameter     = "1.3.6.1.1.1.1.23" // nisSchema.1.23
	AttributeOIDBootFile          = "1.3.6.1.1.1.1.24" // nisSchema.1.24
	AttributeOIDNisMapName        = "1.3.6.1.1.1.1.26" // nisSchema.1.26
	AttributeOIDNisMapEntry       = "1.3.6.1.1.1.1.27" // nisSchema.1.27
)

// NIS Schema Object Class OIDs (RFC 2307).
const (
	ObjectClassOIDPosixAccount   = "1.3.6.1.1.1.2.0"  // nisSchema.2.0
	ObjectClassOIDShadowAccount  = "1.3.6.1.1.1.2.1"  // nisSchema.2.1
	ObjectClassOIDPosixGroup     = "1.3.6.1.1.1.2.2"  // nisSchema.2.2
	ObjectClassOIDIpService      = "1.3.6.1.1.1.2.3"  // nisSchema.2.3
	ObjectClassOIDIpProtocol     = "1.3.6.1.1.1.2.4"  // nisSchema.2.4
	ObjectClassOIDOncRpc         = "1.3.6.1.1.1.2.5"  // nisSchema.2.5
	ObjectClassOIDIpHost         = "1.3.6.1.1.1.2.6"  // nisSchema.2.6
	ObjectClassOIDIpNetwork      = "1.3.6.1.1.1.2.7"  // nisSchema.2.7
	ObjectClassOIDNisNetgroup    = "1.3.6.1.1.1.2.8"  // nisSchema.2.8
	ObjectClassOIDNisMap         = "1.3.6.1.1.1.2.9"  // nisSchema.2.9
	ObjectClassOIDNisObject      = "1.3.6.1.1.1.2.10" // nisSchema.2.10
	ObjectClassOIDIeee802Device  = "1.3.6.1.1.1.2.11" // nisSchema.2.11
	ObjectClassOIDBootableDevice = "1.3.6.1.1.1.2.12" // nisSchema.2.12
)
