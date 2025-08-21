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

// Entry is the struct that represents an LDAP-compatible Entry.
type Entry struct {
	// db data
	Key string    `redis:",key"`  //
	Ver int64     `redis:",ver"`  //
	Ext time.Time `redis:",exat"` //

	// element specific meta data
	Type   attrEntryType   `json:"type,omitempty"`   //  Entry's type `(domain|group|user|host)`
	Status attrEntryStatus `json:"status,omitempty"` //
	BaseDN mod_dn.DN       `json:"baseDN,omitempty"` //

	// element meta data
	UUID            uuid.UUID     `json:"uuid,omitempty"            ldap:"entryUUID"`       //  must be unique
	DN              mod_dn.DN     `json:"dn,omitempty"              ldap:"entryDN"`         //  must be unique
	ObjectClass     []string      `json:"objectClass,omitempty"     ldap:"objectClass"`     //  Entry type
	CreatorsName    mod_dn.DN     `json:"creatorsName,omitempty"    ldap:"creatorsName"`    //
	CreateTimestamp mod_time.Time `json:"createTimestamp,omitempty" ldap:"createTimestamp"` //
	ModifiersName   mod_dn.DN     `json:"modifiersName,omitempty"   ldap:"modifiersName"`   //
	ModifyTimestamp mod_time.Time `json:"modifyTimestamp,omitempty" ldap:"modifyTimestamp"` //

	// element data
	CN                   string         `json:"cn,omitempty"                   ldap:"cn"`                   //  RDN in group's context
	DC                   string         `json:"dc,omitempty"                   ldap:"dc"`                   //
	Description          string         `json:"description,omitempty"          ldap:"description"`          //
	DestinationIndicator []string       `json:"destinationIndicator,omitempty" ldap:"destinationIndicator"` //
	DisplayName          string         `json:"displayName,omitempty"          ldap:"displayName"`          //
	GIDNumber            uint64         `json:"gidNumber,omitempty"            ldap:"gidNumber"`            //  Primary GIDNumber in user's context (ignore it), GIDNumber in group's context.
	HomeDirectory        string         `json:"homeDirectory,omitempty"        ldap:"homeDirectory"`        //
	IPHostNumber         []netip.Prefix `json:"ipHostNumber,omitempty"         ldap:"ipHostNumber"`         //
	Mail                 []string       `json:"mail,omitempty"                 ldap:"mail"`                 //
	Member               []mod_dn.DN    `json:"member,omitempty"               ldap:"member"`               //
	O                    string         `json:"o,omitempty"                    ldap:"o"`                    //
	OU                   string         `json:"ou,omitempty"                   ldap:"ou"`                   //
	Owner                []mod_dn.DN    `json:"owner,omitempty"                ldap:"owner"`                //
	SN                   string         `json:"sn,omitempty"                   ldap:"sn"`                   //
	SSHPublicKey         []string       `json:"sshPublicKey,omitempty"         ldap:"sshPublicKey"`         //
	TelephoneNumber      []string       `json:"telephoneNumber,omitempty"      ldap:"telephoneNumber"`      //
	TelexNumber          []string       `json:"telexNumber,omitempty"          ldap:"telexNumber"`          //
	UID                  string         `json:"uid,omitempty"                  ldap:"uid"`                  //  RDN in user's context
	UIDNumber            uint64         `json:"uidNumber,omitempty"            ldap:"uidNumber"`            //
	UserPassword         string         `json:"userPassword,omitempty"         ldap:"userPassword"`         //
	// UserPKCS12           mod_crypto.Certificates   `json:"userPKCS12,omitempty"           ldap:"userPKCS12"           `           //
	// MemberOf             []mod_dn.DN                   `json:"memberOf,omitempty"             ldap:"memberOf"                         ` //  don't trust LDAP

	// specific data
	AAA string `json:"host_aaa,omitempty"` //  Entry's AAA (?) `(UserPKCS12|UserPassword|SSHPublicKey|etc)`
	ACL string `json:"host_acl,omitempty"` //  Entry's ACL

	// host specific data
	HostType        string     `json:"host_type,omitempty"`         //  host type `(provider|interim|openvpn|ciscovpn)`
	HostASN         uint32     `json:"host_asn,omitempty"`          //
	HostUpstreamASN uint32     `json:"host_upstream_asn,omitempty"` //  upstream route
	HostHostingUUID uuid.UUID  `json:"host_hosting_uuid,omitempty"` //  (?) replace with member/memberOf
	HostURL         *url.URL   `json:"host_url,omitempty"`          //
	HostListen      netip.Addr `json:"host_listen,omitempty"`       //

	// specific data (space-separated KV DB stored as labeledURI)
	LabeledURI []string `json:"labeledURI,omitempty" ldap:"labeledURI"` //
}

// CreateEntryIndex creates the RediSearch index for the Entry struct.
func (r *RedisRepository) CreateEntryIndex() (err error) {
	return r.entry.CreateIndex(r.ctx, func(schema om.FtCreateSchema) rueidis.Completed {
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
