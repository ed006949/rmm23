package mod_db

import (
	"context"
	"net/netip"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/mod_strings"
)

// Entry is the struct that represents an LDAP-compatible Entry.
//
// when updating @src/mod_db/elementEntry.go don't forget to update:
//
//	@src/mod_db/elementRepository.go
type Entry struct {
	// db data
	Key string    `redis:",key"`  //
	Ver int64     `redis:",ver"`  //
	Ext time.Time `redis:",exat"` //

	// element specific meta data
	Type   attrEntryType   `json:"type,omitempty"   msgpack:"type"`   //  Entry's type `(domain|group|user|host)`
	Status attrEntryStatus `json:"status,omitempty" msgpack:"status"` //
	BaseDN attrDN          `json:"baseDN,omitempty" msgpack:"baseDN"` //

	// element meta data
	UUID            uuid.UUID `json:"uuid,omitempty"            ldap:"entryUUID"       msgpack:"uuid"`            //  must be unique
	DN              attrDN    `json:"dn,omitempty"              ldap:"entryDN"         msgpack:"dn"`              //  must be unique
	ObjectClass     []string  `json:"objectClass,omitempty"     ldap:"objectClass"     msgpack:"objectClass"`     //  Entry type
	CreatorsName    attrDN    `json:"creatorsName,omitempty"    ldap:"creatorsName"    msgpack:"creatorsName"`    //
	CreateTimestamp attrTime  `json:"createTimestamp,omitempty" ldap:"createTimestamp" msgpack:"createTimestamp"` //
	ModifiersName   attrDN    `json:"modifiersName,omitempty"   ldap:"modifiersName"   msgpack:"modifiersName"`   //
	ModifyTimestamp attrTime  `json:"modifyTimestamp,omitempty" ldap:"modifyTimestamp" msgpack:"modifyTimestamp"` //

	// element data
	CN                   string         `json:"cn,omitempty"                   ldap:"cn"                   msgpack:"cn"`                   //  RDN in group's context
	DC                   string         `json:"dc,omitempty"                   ldap:"dc"                   msgpack:"dc"`                   //
	Description          string         `json:"description,omitempty"          ldap:"description"          msgpack:"description"`          //
	DestinationIndicator []string       `json:"destinationIndicator,omitempty" ldap:"destinationIndicator" msgpack:"destinationIndicator"` //
	DisplayName          string         `json:"displayName,omitempty"          ldap:"displayName"          msgpack:"displayName"`          //
	GIDNumber            uint64         `json:"gidNumber,omitempty"            ldap:"gidNumber"            msgpack:"gidNumber"`            //  Primary GIDNumber in user's context (ignore it), GIDNumber in group's context.
	HomeDirectory        string         `json:"homeDirectory,omitempty"        ldap:"homeDirectory"        msgpack:"homeDirectory"`        //
	IPHostNumber         []netip.Prefix `json:"ipHostNumber,omitempty"         ldap:"ipHostNumber"         msgpack:"ipHostNumber"`         //
	Mail                 []string       `json:"mail,omitempty"                 ldap:"mail"                 msgpack:"mail"`                 //
	Member               []attrDN       `json:"member,omitempty"               ldap:"member"               msgpack:"member"`               //
	O                    string         `json:"o,omitempty"                    ldap:"o"                    msgpack:"o"`                    //
	OU                   string         `json:"ou,omitempty"                   ldap:"ou"                   msgpack:"ou"`                   //
	Owner                []attrDN       `json:"owner,omitempty"                ldap:"owner"                msgpack:"owner"`                //
	SN                   string         `json:"sn,omitempty"                   ldap:"sn"                   msgpack:"sn"`                   //
	SSHPublicKey         []string       `json:"sshPublicKey,omitempty"         ldap:"sshPublicKey"         msgpack:"sshPublicKey"`         //
	TelephoneNumber      []string       `json:"telephoneNumber,omitempty"      ldap:"telephoneNumber"      msgpack:"telephoneNumber"`      //
	TelexNumber          []string       `json:"telexNumber,omitempty"          ldap:"telexNumber"          msgpack:"telexNumber"`          //
	UID                  string         `json:"uid,omitempty"                  ldap:"uid"                  msgpack:"uid"`                  //  RDN in user's context
	UIDNumber            uint64         `json:"uidNumber,omitempty"            ldap:"uidNumber"            msgpack:"uidNumber"`            //
	UserPassword         string         `json:"userPassword,omitempty"         ldap:"userPassword"         msgpack:"userPassword"`         //
	// UserPKCS12           mod_crypto.Certificates   `json:"userPKCS12,omitempty"           ldap:"userPKCS12"           msgpack:"userPKCS12"`           //
	// MemberOf             []*attrDN                   `json:"memberOf,omitempty"             ldap:"memberOf"             msgpack:"memberOf"            ` //  don't trust LDAP

	// specific data
	AAA string `json:"host_aaa,omitempty" msgpack:"host_aaa"` //  Entry's AAA (?) `(UserPKCS12|UserPassword|SSHPublicKey|etc)`
	ACL string `json:"host_acl,omitempty" msgpack:"host_acl"` //  Entry's ACL

	// host specific data
	HostType        string     `json:"host_type,omitempty"         msgpack:"host_type"`         //  host type `(provider|interim|openvpn|ciscovpn)`
	HostASN         uint32     `json:"host_asn,omitempty"          msgpack:"host_asn"`          //
	HostUpstreamASN uint32     `json:"host_upstream_asn,omitempty" msgpack:"host_upstream_asn"` //  upstream route
	HostHostingUUID uuid.UUID  `json:"host_hosting_uuid,omitempty" msgpack:"host_hosting_uuid"` //  (?) replace with member/memberOf
	HostURL         *url.URL   `json:"host_url,omitempty"          msgpack:"host_url"`          //
	HostListen      netip.Addr `json:"host_listen,omitempty"       msgpack:"host_listen"`       //

	// specific data (space-separated KV DB stored as labeledURI)
	LabeledURI []string `json:"labeledURI,omitempty" ldap:"labeledURI" msgpack:"labeledURI"` //
}

// CreateEntryIndex creates the RediSearch index for the Entry struct.
func (r *RedisRepository) CreateEntryIndex(ctx context.Context) (err error) {
	return r.entry.CreateIndex(ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			FieldName(mod_strings.F_type.FieldName()).As(mod_strings.F_type.String()).Numeric().
			FieldName(mod_strings.F_status.FieldName()).As(mod_strings.F_status.String()).Numeric().
			FieldName(mod_strings.F_baseDN.FieldName()).As(mod_strings.F_baseDN.String()).Tag().Separator(mod_strings.SliceSeparator).

			//
			FieldName(mod_strings.F_uuid.FieldName()).As(mod_strings.F_uuid.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_dn.FieldName()).As(mod_strings.F_dn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_objectClass.FieldNameSlice()).As(mod_strings.F_objectClass.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_creatorsName.FieldName()).As(mod_strings.F_creatorsName.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_createTimestamp.FieldName()).As(mod_strings.F_createTimestamp.String()).Numeric().
			FieldName(mod_strings.F_modifiersName.FieldName()).As(mod_strings.F_modifiersName.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_modifyTimestamp.FieldName()).As(mod_strings.F_modifyTimestamp.String()).Numeric().

			//
			FieldName(mod_strings.F_cn.FieldName()).As(mod_strings.F_cn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_dc.FieldName()).As(mod_strings.F_dc.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_description.FieldName()).As(	mod_strings.F_description.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_destinationIndicator.FieldNameSlice()).As(mod_strings.F_destinationIndicator.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_displayName.FieldName()).As(	mod_strings.F_displayName.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_gidNumber.FieldName()).As(mod_strings.F_gidNumber.String()).Numeric().
			// FieldName(	mod_strings.F_homeDirectory.FieldName()).As(	mod_strings.F_homeDirectory.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_ipHostNumber.FieldNameSlice()).As(mod_strings.F_ipHostNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_mail.FieldNameSlice()).As(mod_strings.F_mail.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_member.FieldNameSlice()).As(mod_strings.F_member.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_o.FieldName()).As(	mod_strings.F_o.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_ou.FieldName()).As(	mod_strings.F_ou.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_owner.FieldNameSlice()).As(mod_strings.F_owner.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_sn.FieldName()).As(	mod_strings.F_sn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_sshPublicKey.FieldNameSlice()).As(mod_strings.F_sshPublicKey.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_telephoneNumber.FieldNameSlice()).As(mod_strings.F_telephoneNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_telexNumber.FieldNameSlice()).As(mod_strings.F_telexNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uid.FieldName()).As(mod_strings.F_uid.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uidNumber.FieldName()).As(mod_strings.F_uidNumber.String()).Numeric().
			// FieldName(mod_strings.F_userPKCS12.FieldNameSlice()).As(mod_strings.F_userPKCS12.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_userPassword.FieldName()).As(	mod_strings.F_userPassword.String()).Tag().Separator(mod_strings.SliceSeparator).

			//
			// FieldName(	mod_strings.F_host_aaa.FieldName()).As(	mod_strings.F_host_aaa.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_acl.FieldName()).As(	mod_strings.F_host_acl.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_type.FieldName()).As(	mod_strings.F_host_type.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_asn.FieldName()).As(	mod_strings.F_host_asn.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_upstream_asn.FieldName()).As(	mod_strings.F_host_upstream_asn.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_hosting_uuid.FieldName()).As(	mod_strings.F_host_hosting_uuid.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_url.FieldName()).As(	mod_strings.F_host_url.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_listen.FieldName()).As(	mod_strings.F_host_listen.String()).Tag().Separator(mod_strings.SliceSeparator).

			//
			// FieldName(mod_strings.F_labeledURI.FieldNameSlice()).As(mod_strings.F_labeledURI.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName( _labeledURI.FieldNameSlice() + ".key").As(	mod_strings.F_labeledURI.String() + "_key").Tag().Separator(mod_strings.SliceSeparator).
			// FieldName( _labeledURI.FieldNameSlice() + ".value").As(	mod_strings.F_labeledURI.String() + "_value").Tag().Separator(mod_strings.SliceSeparator).

			//
			Build()
	})
}
