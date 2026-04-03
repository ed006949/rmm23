package mod_db

import (
	"time"

	"github.com/google/uuid"

	"rmm23/src/mod_dn"
	"rmm23/src/mod_time"
)

type DBEntry struct {
	// db mandatory
	Key string    `redis:",key"`  //
	Ver int64     `redis:",ver"`  //
	Ext time.Time `redis:",exat"` //

	// internal admin (not exposed to LDAP clients)
	Status attrEntryStatus `json:"status,omitempty"` //
	BaseDN mod_dn.DN       `json:"baseDN"`           //  partition key for multi-domain separation

	// objectClass registry (schema authority)
	ObjectClasses         ObjectClassList `json:"objectClasses,omitempty"`                                      //
	objectClassRaw        []string        `json:"-"                               ldap:"objectClass"`           //  unmarshal bridge, not stored in Redis
	StructuralObjectClass string          `json:"structuralObjectClass,omitempty" ldap:"structuralObjectClass"` //  LDAP operational

	// LDAP operational (RFC 4512/4530, exposed to clients, NoUserModify)
	UUID            uuid.UUID     `json:"uuid,omitempty"  ldap:"entryUUID"`       //  must be unique
	DN              mod_dn.DN     `json:"dn"              ldap:"entryDN"`         //  must be unique
	CreatorsName    mod_dn.DN     `json:"creatorsName"    ldap:"creatorsName"`    //
	CreateTimestamp mod_time.Time `json:"createTimestamp" ldap:"createTimestamp"` //
	ModifiersName   mod_dn.DN     `json:"modifiersName"   ldap:"modifiersName"`   //
	ModifyTimestamp mod_time.Time `json:"modifyTimestamp" ldap:"modifyTimestamp"` //
}

// SyncObjectClasses converts the LDAP unmarshal bridge (objectClassRaw) into ObjectClasses.
// Call this after UnmarshalLDAPEntries to populate the schema-aware ObjectClassList.
func (d *DBEntry) SyncObjectClasses() {
	d.ObjectClasses = FromLDAPObjectClass(d.objectClassRaw, nil)
	d.objectClassRaw = nil
}
