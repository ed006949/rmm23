package mod_db

import (
	"net/netip"
	"time"

	"rmm23/src/mod_net"
)

// Entry is the struct that represents an LDAP-compatible Entry.
//
// when updating @src/mod_db/entry_type.go don't forget to update:
//
//	@src/mod_db/entry_const.go
//	@src/mod_db/redis_*.go
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
	UUID            attrUUID          `json:"uuid,omitempty"            ldap:"entryUUID"       msgpack:"uuid"`            //  must be unique
	DN              attrDN            `json:"dn,omitempty"              ldap:"dn"              msgpack:"dn"`              //  must be unique
	ObjectClass     attrObjectClasses `json:"objectClass,omitempty"     ldap:"objectClass"     msgpack:"objectClass"`     //  Entry type
	CreatorsName    attrDN            `json:"creatorsName,omitempty"    ldap:"creatorsName"    msgpack:"creatorsName"`    //
	CreateTimestamp attrTime          `json:"createTimestamp,omitempty" ldap:"createTimestamp" msgpack:"createTimestamp"` //
	ModifiersName   attrDN            `json:"modifiersName,omitempty"   ldap:"modifiersName"   msgpack:"modifiersName"`   //
	ModifyTimestamp attrTime          `json:"modifyTimestamp,omitempty" ldap:"modifyTimestamp" msgpack:"modifyTimestamp"` //

	// element data
	CN                   attrString                `json:"cn,omitempty"                   ldap:"cn"                   msgpack:"cn"`                   //  RDN in group's context
	DC                   attrString                `json:"dc,omitempty"                   ldap:"dc"                   msgpack:"dc"`                   //
	Description          attrString                `json:"description,omitempty"          ldap:"description"          msgpack:"description"`          //
	DestinationIndicator attrDestinationIndicators `json:"destinationIndicator,omitempty" ldap:"destinationIndicator" msgpack:"destinationIndicator"` //
	DisplayName          attrString                `json:"displayName,omitempty"          ldap:"displayName"          msgpack:"displayName"`          //
	GIDNumber            attrIDNumber              `json:"gidNumber,omitempty"            ldap:"gidNumber"            msgpack:"gidNumber"`            //  Primary GIDNumber in user's context (ignore it), GIDNumber in group's context.
	HomeDirectory        attrString                `json:"homeDirectory,omitempty"        ldap:"homeDirectory"        msgpack:"homeDirectory"`        //
	IPHostNumber         attrIPHostNumbers         `json:"ipHostNumber,omitempty"         ldap:"ipHostNumber"         msgpack:"ipHostNumber"`         //
	Mail                 attrMails                 `json:"mail,omitempty"                 ldap:"mail"                 msgpack:"mail"`                 //
	Member               attrDNs                   `json:"member,omitempty"               ldap:"member"               msgpack:"member"`               //
	O                    attrString                `json:"o,omitempty"                    ldap:"o"                    msgpack:"o"`                    //
	OU                   attrString                `json:"ou,omitempty"                   ldap:"ou"                   msgpack:"ou"`                   //
	Owner                attrDNs                   `json:"owner,omitempty"                ldap:"owner"                msgpack:"owner"`                //
	SN                   attrString                `json:"sn,omitempty"                   ldap:"sn"                   msgpack:"sn"`                   //
	SSHPublicKey         attrSSHPublicKeys         `json:"sshPublicKey,omitempty"         ldap:"sshPublicKey"         msgpack:"sshPublicKey"`         //
	TelephoneNumber      attrStrings               `json:"telephoneNumber,omitempty"      ldap:"telephoneNumber"      msgpack:"telephoneNumber"`      //
	TelexNumber          attrStrings               `json:"telexNumber,omitempty"          ldap:"telexNumber"          msgpack:"telexNumber"`          //
	UID                  attrID                    `json:"uid,omitempty"                  ldap:"uid"                  msgpack:"uid"`                  //  RDN in user's context
	UIDNumber            attrIDNumber              `json:"uidNumber,omitempty"            ldap:"uidNumber"            msgpack:"uidNumber"`            //
	// UserPKCS12           mod_crypto.Certificates   `json:"userPKCS12,omitempty"           ldap:"userPKCS12"           msgpack:"userPKCS12"`           //
	UserPassword attrUserPassword `json:"userPassword,omitempty" ldap:"userPassword" msgpack:"userPassword"` //
	// MemberOf             attrDNs                   `json:"memberOf,omitempty"             ldap:"memberOf"             msgpack:"memberOf"            ` //  don't trust LDAP

	// specific data
	AAA string `json:"host_aaa,omitempty" msgpack:"host_aaa"` //  Entry's AAA (?) `(UserPKCS12|UserPassword|SSHPublicKey|etc)`
	ACL string `json:"host_acl,omitempty" msgpack:"host_acl"` //  Entry's ACL

	// host specific data
	HostType        string       `json:"host_type,omitempty"         msgpack:"host_type"`         //  host type `(provider|interim|openvpn|ciscovpn)`
	HostASN         uint32       `json:"host_asn,omitempty"          msgpack:"host_asn"`          //
	HostUpstreamASN uint32       `json:"host_upstream_asn,omitempty" msgpack:"host_upstream_asn"` //  upstream route
	HostHostingUUID uint32       `json:"host_hosting_uuid,omitempty" msgpack:"host_hosting_uuid"` //  (?) replace with member/memberOf
	HostURL         *mod_net.URL `json:"host_url,omitempty"          msgpack:"host_url"`          //
	HostListen      *netip.Addr  `json:"host_listen,omitempty"       msgpack:"host_listen"`       //

	// specific data (space-separated KV DB stored as labeledURI)
	LabeledURI attrLabeledURIs `json:"labeledURI,omitempty" ldap:"labeledURI" msgpack:"labeledURI"` //
}
