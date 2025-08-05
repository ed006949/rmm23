package mod_db

import (
	"math/big"
	"net/netip"
	"net/url"
	"time"

	"rmm23/src/mod_bools"
	"rmm23/src/mod_crypto"
)

// Cert is the struct that represents an LDAP userPKCS12 attribute.
//
// when updating @src/mod_db/entry_type.go don't forget to update:
//
//	@src/mod_db/certificate_*.go
//	@src/mod_db/redis_field.go
//	@src/mod_db/redis_*.go
type Cert struct {
	// db data
	Key string    `redis:",key"`  //
	Ver int64     `redis:",ver"`  //
	Ext time.Time `redis:",exat"` //

	// element meta data
	UUID           attrUUID           `json:"uuid"           msgpack:"uuid"`           // x509.Certificate.Raw() hash
	SerialNumber   *big.Int           `json:"serialNumber"   msgpack:"serialNumber"`   // (?) redis:",key"
	Issuer         attrDN             `json:"issuer"         msgpack:"issuer"`         //
	Subject        attrDN             `json:"subject"        msgpack:"subject"`        //
	NotBefore      attrTime           `json:"notBefore"      msgpack:"notBefore"`      //
	NotAfter       attrTime           `json:"notAfter"       msgpack:"notAfter"`       // redis:",exat"
	DNSNames       []string           `json:"dnsNames"       msgpack:"dnsNames"`       //
	EmailAddresses []string           `json:"emailAddresses" msgpack:"emailAddresses"` //
	IPAddresses    []*netip.Addr      `json:"ipAddresses"    msgpack:"ipAddresses"`    //
	URIs           []*url.URL         `json:"uris"           msgpack:"uris"`           //
	IsCA           mod_bools.AttrBool `json:"isCA"           msgpack:"isCA"`           //

	// // element specific meta data
	// Type   attrEntryType   `json:"type,omitempty"   msgpack:"type"`   // (?) Certificate's type
	// Status attrEntryStatus `json:"status,omitempty" msgpack:"status"` //
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

}
