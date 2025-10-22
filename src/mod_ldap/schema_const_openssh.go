package mod_ldap

// OpenSSH LPK (LDAP Public Key) Schema OIDs
// Based on common OpenSSH LDAP implementations.
const (
	// OpenSSHLPKBase is the OpenSSH LPK base OID (commonly used).
	OpenSSHLPKBase = "1.3.6.1.4.1.24552.500.1.1"
)

// OpenSSH LPK Attribute OIDs.
const (
	AttributeOIDSshPublicKey = "1.3.6.1.4.1.24552.500.1.1.1.13"
)

// OpenSSH LPK Object Class OIDs.
const (
	ObjectClassOIDLdapPublicKey = "1.3.6.1.4.1.24552.500.1.1.2.0"
)
