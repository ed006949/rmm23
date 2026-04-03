package mod_db

import (
	"math/big"
	"net/netip"
	"net/url"

	"github.com/google/uuid"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/mod_crypto"
	"rmm23/src/mod_dn"
	"rmm23/src/mod_strings"
	"rmm23/src/mod_time"
)

// Entry is the struct that represents an LDAP-compatible Entry.
type Entry struct {
	DBEntry

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
	HostListen      netip.Addr `json:"host_listen"`                 //

	// specific data (space-separated KV DB stored as labeledURI)
	LabeledURI []string `json:"labeledURI,omitempty" ldap:"labeledURI"` //

	// certificate attributes (objectClass: pkiUser / pkiCA)
	SerialNumber   *big.Int                `json:"serialNumber,omitempty"`   //
	Issuer         mod_dn.DN               `json:"issuer"`                   //
	Subject        mod_dn.DN               `json:"subject"`                  //
	NotBefore      mod_time.Time           `json:"notBefore"`                //
	NotAfter       mod_time.Time           `json:"notAfter"`                 //
	DNSNames       []string                `json:"dnsNames,omitempty"`       //
	EmailAddresses []string                `json:"emailAddresses,omitempty"` //
	IPAddresses    []*netip.Addr           `json:"ipAddresses,omitempty"`    //
	URIs           []*url.URL              `json:"uris,omitempty"`           //
	IsCA           bool                    `json:"isCA,omitempty"`           //
	Certificate    *mod_crypto.Certificate `json:"certificate,omitempty"`    //
}

// CreateEntryIndex creates the RediSearch index for the Entry struct.
func (r *RedisRepository) CreateEntryIndex() (err error) {
	return r.entry.CreateIndex(r.ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			// internal admin
			FieldName(mod_strings.F_status.FieldName()).As(mod_strings.F_status.String()).Numeric().
			FieldName(mod_strings.F_baseDN.FieldName()).As(mod_strings.F_baseDN.String()).Tag().Separator(mod_strings.SliceSeparator).

			// objectClass (JSON path into ObjectClassList)
			FieldName("$.objectClasses[*].name").As(mod_strings.F_objectClass.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_structuralObjectClass.FieldName()).As(mod_strings.F_structuralObjectClass.String()).Tag().Separator(mod_strings.SliceSeparator).

			// LDAP operational
			FieldName(mod_strings.F_uuid.FieldName()).As(mod_strings.F_uuid.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_dn.FieldName()).As(mod_strings.F_dn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_creatorsName.FieldName()).As(mod_strings.F_creatorsName.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_createTimestamp.FieldName()).As(mod_strings.F_createTimestamp.String()).Numeric().
			FieldName(mod_strings.F_modifiersName.FieldName()).As(mod_strings.F_modifiersName.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_modifyTimestamp.FieldName()).As(mod_strings.F_modifyTimestamp.String()).Numeric().

			// standard LDAP attributes
			FieldName(mod_strings.F_cn.FieldName()).As(mod_strings.F_cn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_dc.FieldName()).As(mod_strings.F_dc.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(mod_strings.F_description.FieldName()).As(mod_strings.F_description.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_destinationIndicator.FieldNameSlice()).As(mod_strings.F_destinationIndicator.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(mod_strings.F_displayName.FieldName()).As(mod_strings.F_displayName.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_gidNumber.FieldName()).As(mod_strings.F_gidNumber.String()).Numeric().
			// FieldName(mod_strings.F_homeDirectory.FieldName()).As(mod_strings.F_homeDirectory.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_ipHostNumber.FieldNameSlice()).As(mod_strings.F_ipHostNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_mail.FieldNameSlice()).As(mod_strings.F_mail.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_member.FieldNameSlice()).As(mod_strings.F_member.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(mod_strings.F_o.FieldName()).As(mod_strings.F_o.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(mod_strings.F_ou.FieldName()).As(mod_strings.F_ou.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_owner.FieldNameSlice()).As(mod_strings.F_owner.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(mod_strings.F_sn.FieldName()).As(mod_strings.F_sn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_sshPublicKey.FieldNameSlice()).As(mod_strings.F_sshPublicKey.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_telephoneNumber.FieldNameSlice()).As(mod_strings.F_telephoneNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_telexNumber.FieldNameSlice()).As(mod_strings.F_telexNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uid.FieldName()).As(mod_strings.F_uid.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uidNumber.FieldName()).As(mod_strings.F_uidNumber.String()).Numeric().
			// FieldName(mod_strings.F_userPassword.FieldName()).As(mod_strings.F_userPassword.String()).Tag().Separator(mod_strings.SliceSeparator).

			// certificate attributes (objectClass: pkiUser / pkiCA)
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
