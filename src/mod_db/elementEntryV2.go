package mod_db

import (
	"net/netip"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/mod_dn"
	"rmm23/src/mod_strings"
	"rmm23/src/mod_time"
)

// EntryV2 is an enhanced LDAP-compatible Entry with objectClass sub-attribute support
// This struct supports:
// 1. Redis-specific attributes (Key, Ver, ExAt)
// 2. LDAP operational attributes (entryDN, entryUUID, timestamps, etc.)
// 3. objectClass with schema-specific sub-attributes (2307bis, etc.)
// 4. RediSearch-optimized indexing for all attribute types.
type EntryV2 struct {
	// ========================================
	// Redis-specific attributes
	// ========================================
	Key  string    `redis:",key"`  // Redis key
	Ver  int64     `redis:",ver"`  // Version for optimistic locking
	ExAt time.Time `redis:",exat"` // Expiration time

	// ========================================
	// LDAP Operational Attributes (RFC 4512)
	// ========================================
	// Distinguished Name and UUID
	EntryDN   mod_dn.DN `json:"entryDN,omitempty"   ldap:"entryDN"`   // DN of the entry
	EntryUUID uuid.UUID `json:"entryUUID,omitempty" ldap:"entryUUID"` // UUID of the entry

	// Timestamps and creators (operational attributes)
	CreateTimestamp mod_time.Time `json:"createTimestamp,omitempty" ldap:"createTimestamp"` // Creation time
	ModifyTimestamp mod_time.Time `json:"modifyTimestamp,omitempty" ldap:"modifyTimestamp"` // Last modification time
	CreatorsName    mod_dn.DN     `json:"creatorsName,omitempty"    ldap:"creatorsName"`    // Creator DN
	ModifiersName   mod_dn.DN     `json:"modifiersName,omitempty"   ldap:"modifiersName"`   // Last modifier DN

	// Structural information (operational attributes)
	StructuralObjectClass string    `json:"structuralObjectClass,omitempty" ldap:"structuralObjectClass"` // Structural objectClass
	SubschemaSubentry     mod_dn.DN `json:"subschemaSubentry,omitempty"     ldap:"subschemaSubentry"`     // Schema entry DN
	HasSubordinates       bool      `json:"hasSubordinates,omitempty"       ldap:"hasSubordinates"`       // Has children
	SubordinateCount      int64     `json:"subordinateCount,omitempty"      ldap:"subordinateCount"`      // Number of children

	// Replication attributes (operational)
	EntryCSN   string `json:"entryCSN,omitempty"   ldap:"entryCSN"`   // Change sequence number
	ContextCSN string `json:"contextCSN,omitempty" ldap:"contextCSN"` // Context CSN

	// ========================================
	// objectClass with Sub-Attributes Support
	// ========================================
	// ObjectClassData stores objectClass names and their schema-specific attributes
	// This allows storing 2307bis attributes, custom schema attributes, etc.
	ObjectClassData ObjectClassList `json:"objectClassData,omitempty"` // objectClass with sub-attributes

	// ObjectClass is a computed field for backward compatibility and TAG search
	// It's derived from ObjectClassData.Names()
	ObjectClass []string `json:"objectClass,omitempty" ldap:"objectClass"` // objectClass names for search

	// ========================================
	// Standard LDAP Attributes
	// ========================================
	// Core attributes
	CN          string `json:"cn,omitempty"          ldap:"cn"`          // Common Name
	DC          string `json:"dc,omitempty"          ldap:"dc"`          // Domain Component
	Description string `json:"description,omitempty" ldap:"description"` // Description
	DisplayName string `json:"displayName,omitempty" ldap:"displayName"` // Display Name
	O           string `json:"o,omitempty"           ldap:"o"`           // Organization
	OU          string `json:"ou,omitempty"          ldap:"ou"`          // Organizational Unit
	SN          string `json:"sn,omitempty"          ldap:"sn"`          // Surname

	// POSIX attributes (RFC 2307bis)
	UID           string `json:"uid,omitempty"           ldap:"uid"`           // User ID
	UIDNumber     uint64 `json:"uidNumber,omitempty"     ldap:"uidNumber"`     // Numeric UID
	GIDNumber     uint64 `json:"gidNumber,omitempty"     ldap:"gidNumber"`     // Numeric GID
	HomeDirectory string `json:"homeDirectory,omitempty" ldap:"homeDirectory"` // Home directory

	// Group attributes
	Member []mod_dn.DN `json:"member,omitempty" ldap:"member"` // Group members
	Owner  []mod_dn.DN `json:"owner,omitempty"  ldap:"owner"`  // Group owners

	// Network attributes
	IPHostNumber []netip.Prefix `json:"ipHostNumber,omitempty" ldap:"ipHostNumber"` // IP addresses
	Mail         []string       `json:"mail,omitempty"         ldap:"mail"`         // Email addresses

	// Communication attributes
	TelephoneNumber      []string `json:"telephoneNumber,omitempty"      ldap:"telephoneNumber"`      // Phone numbers
	TelexNumber          []string `json:"telexNumber,omitempty"          ldap:"telexNumber"`          // Telex numbers
	DestinationIndicator []string `json:"destinationIndicator,omitempty" ldap:"destinationIndicator"` // Routing indicators

	// Security attributes
	UserPassword string   `json:"userPassword,omitempty" ldap:"userPassword"` // Password (hashed)
	SSHPublicKey []string `json:"sshPublicKey,omitempty" ldap:"sshPublicKey"` // SSH public keys

	// URI attributes
	LabeledURI []string `json:"labeledURI,omitempty" ldap:"labeledURI"` // Labeled URIs (KV storage)

	// ========================================
	// Custom/Application-Specific Attributes
	// ========================================
	// Entry metadata
	Type   attrEntryType   `json:"type,omitempty"`   // Entry type (domain|group|user|host)
	Status attrEntryStatus `json:"status,omitempty"` // Entry status
	BaseDN mod_dn.DN       `json:"baseDN,omitempty"` // Base DN

	// Authentication/Authorization
	AAA string `json:"host_aaa,omitempty"` // Authentication method
	ACL string `json:"host_acl,omitempty"` // Access control list

	// Host-specific attributes
	HostType        string     `json:"host_type,omitempty"`         // Host type (provider|interim|openvpn|ciscovpn)
	HostASN         uint32     `json:"host_asn,omitempty"`          // AS Number
	HostUpstreamASN uint32     `json:"host_upstream_asn,omitempty"` // Upstream AS Number
	HostHostingUUID uuid.UUID  `json:"host_hosting_uuid,omitempty"` // Hosting UUID
	HostURL         *url.URL   `json:"host_url,omitempty"`          // Host URL
	HostListen      netip.Addr `json:"host_listen,omitempty"`       // Listen address
}

// CreateEntryV2Index creates the RediSearch index for EntryV2 with full attribute support.
func (r *RedisRepository) CreateEntryV2Index() (err error) {
	// You would need to add entryV2 repository to RedisRepository first
	// This is a template showing how to index all attributes including objectClass sub-attributes
	return r.entry.CreateIndex(r.ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			// ========================================
			// Custom metadata fields
			// ========================================
			FieldName(mod_strings.F_type.FieldName()).As(mod_strings.F_type.String()).Numeric().
			FieldName(mod_strings.F_status.FieldName()).As(mod_strings.F_status.String()).Numeric().
			FieldName(mod_strings.F_baseDN.FieldName()).As(mod_strings.F_baseDN.String()).Tag().Separator(mod_strings.SliceSeparator).

			// ========================================
			// LDAP Operational Attributes (Searchable)
			// ========================================
			FieldName("$.entryDN").As("entryDN").Tag().Separator(mod_strings.SliceSeparator).
			FieldName("$.entryUUID").As("entryUUID").Tag().Separator(mod_strings.SliceSeparator).
			FieldName("$.createTimestamp").As("createTimestamp").Numeric().
			FieldName("$.modifyTimestamp").As("modifyTimestamp").Numeric().
			FieldName("$.creatorsName").As("creatorsName").Tag().Separator(mod_strings.SliceSeparator).
			FieldName("$.modifiersName").As("modifiersName").Tag().Separator(mod_strings.SliceSeparator).
			FieldName("$.structuralObjectClass").As("structuralObjectClass").Tag().Separator(mod_strings.SliceSeparator).
			FieldName("$.subschemaSubentry").As("subschemaSubentry").Tag().Separator(mod_strings.SliceSeparator).
			FieldName("$.hasSubordinates").As("hasSubordinates").Tag().
			FieldName("$.subordinateCount").As("subordinateCount").Numeric().
			FieldName("$.entryCSN").As("entryCSN").Tag().
			FieldName("$.contextCSN").As("contextCSN").Tag().

			// ========================================
			// objectClass (TAG searchable)
			// ========================================
			FieldName(mod_strings.F_objectClass.FieldNameSlice()).As(mod_strings.F_objectClass.String()).Tag().Separator(mod_strings.SliceSeparator).

			// objectClass sub-attributes as JSON path queries
			// Example: Search for posixAccount.uidNumber or inetOrgPerson.mail
			FieldName("$.objectClassData[*].name").As("objectClassName").Tag().Separator(mod_strings.SliceSeparator).
			// Note: For deep attribute search, you'd query: @objectClassName:{posixAccount}
			// Then fetch and filter by objectClassData attributes in application layer

			// ========================================
			// Standard LDAP Attributes
			// ========================================
			FieldName(mod_strings.F_cn.FieldName()).As(mod_strings.F_cn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_dc.FieldName()).As(mod_strings.F_dc.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uid.FieldName()).As(mod_strings.F_uid.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uidNumber.FieldName()).As(mod_strings.F_uidNumber.String()).Numeric().
			FieldName(mod_strings.F_gidNumber.FieldName()).As(mod_strings.F_gidNumber.String()).Numeric().
			FieldName(mod_strings.F_member.FieldNameSlice()).As(mod_strings.F_member.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_owner.FieldNameSlice()).As(mod_strings.F_owner.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_ipHostNumber.FieldNameSlice()).As(mod_strings.F_ipHostNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_mail.FieldNameSlice()).As(mod_strings.F_mail.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_sshPublicKey.FieldNameSlice()).As(mod_strings.F_sshPublicKey.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_telephoneNumber.FieldNameSlice()).As(mod_strings.F_telephoneNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_telexNumber.FieldNameSlice()).As(mod_strings.F_telexNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_destinationIndicator.FieldNameSlice()).As(mod_strings.F_destinationIndicator.String()).Tag().Separator(mod_strings.SliceSeparator).
			Build()
	})
}

// SyncObjectClass synchronizes ObjectClass slice with ObjectClassData
// Call this before saving to Redis to ensure consistency.
func (e *EntryV2) SyncObjectClass() {
	e.ObjectClass = e.ObjectClassData.Names()
}

// AddObjectClass adds a new objectClass with optional attributes.
func (e *EntryV2) AddObjectClass(className string, attrs map[string]interface{}) {
	if e.ObjectClassData == nil {
		e.ObjectClassData = make(ObjectClassList, 0)
	}

	entry := ObjectClassEntry{
		Name:       className,
		Attributes: attrs,
	}
	if entry.Attributes == nil {
		entry.Attributes = make(map[string]interface{})
	}

	e.ObjectClassData = append(e.ObjectClassData, entry)
	e.SyncObjectClass()
}

// GetObjectClassAttribute retrieves an attribute from a specific objectClass.
func (e *EntryV2) GetObjectClassAttribute(className, attrName string) (interface{}, bool) {
	return e.ObjectClassData.GetAttribute(className, attrName)
}

// SetObjectClassAttribute sets an attribute in a specific objectClass.
func (e *EntryV2) SetObjectClassAttribute(className, attrName string, value interface{}) {
	if oc := e.ObjectClassData.GetClass(className); oc != nil {
		oc.Attributes[attrName] = value
	}
}
