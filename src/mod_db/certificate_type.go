package mod_db

import (
	"time"

	"rmm23/src/mod_crypto"
)

// Cert is the struct that represents an LDAP userPKCS12 attribute.
//
// when updating @src/mod_db/entry_type.go don't forget to update:
//
//	@src/mod_db/certificate_*.go
//	@src/mod_db/redis_*.go
type Cert struct {
	// db data
	Key string    `redis:",key"`  //
	Ver int64     `redis:",ver"`  //
	Ext time.Time `redis:",exat"` //

	// // element specific meta data
	// Type   attrEntryType   `json:"type,omitempty"   msgpack:"type"`   // (?) Certificate's type
	Status attrEntryStatus `json:"status,omitempty" msgpack:"status"` //
	// BaseDN attrDN          `json:"baseDN,omitempty" msgpack:"baseDN"` //

	// // element meta data
	// UUID            attrUUID `json:"uuid,omitempty"            msgpack:"uuid"`            //  must be unique
	// DN              attrDN   `json:"dn,omitempty"              msgpack:"dn"`              //  must be unique
	// CreatorsName    attrDN   `json:"creatorsName,omitempty"    msgpack:"creatorsName"`    //
	// CreateTimestamp attrTime `json:"createTimestamp,omitempty" msgpack:"createTimestamp"` //
	// ModifiersName   attrDN   `json:"modifiersName,omitempty"   msgpack:"modifiersName"`   //
	// ModifyTimestamp attrTime `json:"modifyTimestamp,omitempty" msgpack:"modifyTimestamp"` //

	// element data
	Certificate *mod_crypto.Certificate `json:"certificate,omitempty" msgpack:"certificate"` //

	// // element data
	// CN           attrString        `json:"cn,omitempty"           msgpack:"cn"`           //  RDN in group's context
	// DC           attrString        `json:"dc,omitempty"           msgpack:"dc"`           //
	// Description  attrString        `json:"description,omitempty"  msgpack:"description"`  //
	// IPHostNumber attrIPHostNumbers `json:"ipHostNumber,omitempty" msgpack:"ipHostNumber"` //
	// Mail         attrMails         `json:"mail,omitempty"         msgpack:"mail"`         //
	// O            attrString        `json:"o,omitempty"            msgpack:"o"`            //
	// OU           attrString        `json:"ou,omitempty"           msgpack:"ou"`           //
	// SN           attrString        `json:"sn,omitempty"           msgpack:"sn"`           //
}
