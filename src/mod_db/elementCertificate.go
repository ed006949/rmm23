package mod_db

import (
	"math/big"
	"net/netip"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/mod_crypto"
	"rmm23/src/mod_strings"
)

// Cert is the struct that represents an LDAP userPKCS12 attribute.
type Cert struct {
	// db data
	Key string    `redis:",key"`  //
	Ver int64     `redis:",ver"`  //
	Ext time.Time `redis:",exat"` //

	// element specific meta data
	// Type   attrEntryType   `json:"type,omitempty"   msgpack:"type"`   //
	Status attrEntryStatus `json:"status,omitempty" msgpack:"status"` //
	// BaseDN attrDN          `json:"baseDN,omitempty" msgpack:"baseDN"` //

	// // element meta data
	// UUID            attrUUID `json:"uuid,omitempty"            msgpack:"uuid"`            //  must be unique
	// DN              attrDN   `json:"dn,omitempty"              msgpack:"dn"`              //  must be unique
	// CreatorsName    attrDN   `json:"creatorsName,omitempty"    msgpack:"creatorsName"`    //
	// CreateTimestamp attrTime `json:"createTimestamp,omitempty" msgpack:"createTimestamp"` //
	// ModifiersName   attrDN   `json:"modifiersName,omitempty"   msgpack:"modifiersName"`   //
	// ModifyTimestamp attrTime `json:"modifyTimestamp,omitempty" msgpack:"modifyTimestamp"` //

	// element meta data
	UUID           uuid.UUID     `json:"uuid"           msgpack:"uuid"`           // x509.Certificate.Raw() hash `redis:",key"`
	SerialNumber   *big.Int      `json:"serialNumber"   msgpack:"serialNumber"`   // (?) redis:",key". it can be non-uniq like LDAP's entryUUID - not trusted.
	Issuer         attrDN        `json:"issuer"         msgpack:"issuer"`         //
	Subject        attrDN        `json:"subject"        msgpack:"subject"`        //
	NotBefore      attrTime      `json:"notBefore"      msgpack:"notBefore"`      //
	NotAfter       attrTime      `json:"notAfter"       msgpack:"notAfter"`       // (?) redis:",exat"
	DNSNames       []string      `json:"dnsNames"       msgpack:"dnsNames"`       //
	EmailAddresses []string      `json:"emailAddresses" msgpack:"emailAddresses"` //
	IPAddresses    []*netip.Addr `json:"ipAddresses"    msgpack:"ipAddresses"`    //
	URIs           []*url.URL    `json:"uris"           msgpack:"uris"`           //
	IsCA           bool          `json:"isCA"           msgpack:"isCA"`           //

	// element data
	Certificate *mod_crypto.Certificate `json:"certificate,omitempty" msgpack:"certificate"` //
}

// CreateCertIndex creates the RediSearch index for the Cert struct.
func (r *RedisRepository) CreateCertIndex() (err error) {
	return r.cert.CreateIndex(r.ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.

			//
			// FieldName(mod_strings.F_type.FieldName()).As(mod_strings.F_type.String()).Numeric().
			FieldName(mod_strings.F_status.FieldName()).As(mod_strings.F_status.String()).Numeric().
			// FieldName(mod_strings.F_baseDN.FieldName()).As(mod_strings.F_baseDN.String()).Tag().Separator(mod_strings.SliceSeparator).

			//
			FieldName(mod_strings.F_uuid.FieldName()).As(mod_strings.F_uuid.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_serialNumber.FieldName()).As(mod_strings.F_serialNumber.String()).Numeric().
			FieldName(mod_strings.F_issuer.FieldName()).As(mod_strings.F_issuer.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_subject.FieldName()).As(mod_strings.F_subject.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_notBefore.FieldName()).As(mod_strings.F_notBefore.String()).Numeric().
			FieldName(mod_strings.F_notAfter.FieldName()).As(mod_strings.F_notAfter.String()).Numeric().
			FieldName(mod_strings.F_dnsNames.FieldNameSlice()).As(mod_strings.F_dnsNames.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_emailAddresses.FieldNameSlice()).As(mod_strings.F_emailAddresses.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_ipAddresses.FieldNameSlice()).As(mod_strings.F_ipAddresses.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uris.FieldNameSlice()).As(mod_strings.F_uris.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_isCA.FieldName()).As(mod_strings.F_isCA.String()).Tag().Separator(mod_strings.SliceSeparator).

			//
			Build()
	})
}
