package mod_db

import (
	"net/netip"

	"github.com/RediSearch/redisearch-go/redisearch"

	"rmm23/src/mod_net"
)

type Conf struct {
	URL       *mod_net.URL `json:"url,omitempty"`
	Name      string       `json:"name,omitempty"`
	rsClient  *redisearch.Client
	rcNetwork string
}

// Entry is the struct that represents an LDAP-compatible entry.
type Entry struct {
	// element specific meta data
	Type   AttrEntryType   `json:"type,omitempty"   msgpack:"type,omitempty"   redis:"type"   redisearch:"text,sortable"` // entry's type `(domain|group|user|host)`
	Status AttrEntryStatus `json:"status,omitempty" msgpack:"status,omitempty" redis:"status" redisearch:"text,sortable"` //
	BaseDN attrDN          `json:"baseDN,omitempty" msgpack:"baseDN,omitempty" redis:"baseDN" redisearch:"text,sortable"` //

	// element meta data
	UUID            attrUUID          `json:"uuid,omitempty"            ldap:"entryUUID"       msgpack:"uuid,omitempty"            redis:"uuid"            redisearch:"text,sortable"` // must be unique
	DN              attrDN            `json:"dn,omitempty"              ldap:"dn"              msgpack:"dn,omitempty"              redis:"dn"              redisearch:"text,sortable"` // must be unique
	ObjectClass     attrObjectClasses `json:"objectClass,omitempty"     ldap:"objectClass"     msgpack:"objectClass,omitempty"     redis:"objectClass"     redisearch:"tag"`           // entry type
	CreatorsName    attrDN            `json:"creatorsName,omitempty"    ldap:"creatorsName"    msgpack:"creatorsName,omitempty"    redis:"creatorsName"    redisearch:"text"`          //
	CreateTimestamp attrTimestamp     `json:"createTimestamp,omitempty" ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"`          //
	ModifiersName   attrDN            `json:"modifiersName,omitempty"   ldap:"modifiersName"   msgpack:"modifiersName,omitempty"   redis:"modifiersName"   redisearch:"text"`          //
	ModifyTimestamp attrTimestamp     `json:"modifyTimestamp,omitempty" ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"`          //

	// element data
	CN                   attrString                `json:"cn,omitempty"                   ldap:"cn"                   msgpack:"cn,omitempty"                   redis:"cn"                   redisearch:"text"`             // RDN in group's context
	DC                   attrString                `json:"dc,omitempty"                   ldap:"dc"                   msgpack:"dc,omitempty"                   redis:"dc"                   redisearch:"text,sortable"`    //
	Description          attrString                `json:"description,omitempty"          ldap:"description"          msgpack:"description,omitempty"          redis:"description"          redisearch:"text"`             //
	DestinationIndicator attrDestinationIndicators `json:"destinationIndicator,omitempty" ldap:"destinationIndicator" msgpack:"destinationIndicator,omitempty" redis:"destinationIndicator" redisearch:"text"`             //
	DisplayName          attrString                `json:"displayName,omitempty"          ldap:"displayName"          msgpack:"displayName,omitempty"          redis:"displayName"          redisearch:"text,sortable"`    //
	GIDNumber            attrIDNumber              `json:"gidNumber,omitempty"            ldap:"gidNumber"            msgpack:"gidNumber,omitempty"            redis:"gidNumber"            redisearch:"numeric"`          // Primary GIDNumber in user's context (ignore it) and GIDNumber in group's context.
	HomeDirectory        attrString                `json:"homeDirectory,omitempty"        ldap:"homeDirectory"        msgpack:"homeDirectory,omitempty"        redis:"homeDirectory"        redisearch:"text"`             //
	IPHostNumber         attrIPHostNumbers         `json:"ipHostNumber,omitempty"         ldap:"ipHostNumber"         msgpack:"ipHostNumber,omitempty"         redis:"ipHostNumber"         redisearch:"text,sortable"`    //
	Mail                 attrMails                 `json:"mail,omitempty"                 ldap:"mail"                 msgpack:"mail,omitempty"                 redis:"mail"                 redisearch:"tag"`              //
	Member               attrDNs                   `json:"member,omitempty"               ldap:"member"               msgpack:"member,omitempty"               redis:"member"               redisearch:"tag,sortable"`     //
	O                    attrString                `json:"o,omitempty"                    ldap:"o"                    msgpack:"o,omitempty"                    redis:"o"                    redisearch:"text"`             //
	OU                   attrString                `json:"ou,omitempty"                   ldap:"ou"                   msgpack:"ou,omitempty"                   redis:"ou"                   redisearch:"text"`             //
	Owner                attrDNs                   `json:"owner,omitempty"                ldap:"owner"                msgpack:"owner,omitempty"                redis:"owner"                redisearch:"tag"`              //
	SN                   attrString                `json:"sn,omitempty"                   ldap:"sn"                   msgpack:"sn,omitempty"                   redis:"sn"                   redisearch:"text"`             //
	SSHPublicKey         attrSSHPublicKeys         `json:"sshPublicKey,omitempty"         ldap:"sshPublicKey"         msgpack:"sshPublicKey,omitempty"         redis:"sshPublicKey"         redisearch:"tag"`              //
	TelephoneNumber      attrStrings               `json:"telephoneNumber,omitempty"      ldap:"telephoneNumber"      msgpack:"telephoneNumber,omitempty"      redis:"telephoneNumber"      redisearch:"tag"`              //
	TelexNumber          attrStrings               `json:"telexNumber,omitempty"          ldap:"telexNumber"          msgpack:"telexNumber,omitempty"          redis:"telexNumber"          redisearch:"tag"`              //
	UID                  attrID                    `json:"uid,omitempty"                  ldap:"uid"                  msgpack:"uid,omitempty"                  redis:"uid"                  redisearch:"text,sortable"`    // RDN in user's context
	UIDNumber            attrIDNumber              `json:"uidNumber,omitempty"            ldap:"uidNumber"            msgpack:"uidNumber,omitempty"            redis:"uidNumber"            redisearch:"numeric,sortable"` //
	UserPKCS12           attrUserPKCS12s           `json:"userPKCS12,omitempty"           ldap:"userPKCS12"           msgpack:"userPKCS12,omitempty"           redis:"userPKCS12"           redisearch:"tag"`              //
	UserPassword         attrUserPassword          `json:"userPassword,omitempty"         ldap:"userPassword"         msgpack:"userPassword,omitempty"         redis:"userPassword"         redisearch:"text"`             //
	// MemberOf             attrDNs                   `json:"memberOf,omitempty"             ldap:"memberOf"             msgpack:"memberOf,omitempty"             redis:"memberOf"             redisearch:"tag"`              // ignore it, don't cache, calculate on the fly or avoid

	// specific data
	AAA string `json:"host_aaa,omitempty" msgpack:"host_aaa,omitempty" redis:"host_aaa" redisearch:"text"` // entry's AAA (?) `(UserPKCS12|UserPassword|SSHPublicKey|etc)`
	ACL string `json:"host_acl,omitempty" msgpack:"host_acl,omitempty" redis:"host_acl" redisearch:"text"` // entry's ACL

	// host specific data
	HostType        string       `json:"host_type,omitempty"         msgpack:"host_type,omitempty"         redis:"host_type"         redisearch:"text"`             // host type `(provider|interim|openvpn|ciscovpn)`
	HostASN         uint32       `json:"host_asn,omitempty"          msgpack:"host_asn,omitempty"          redis:"host_asn"          redisearch:"numeric,sortable"` //
	HostUpstreamASN uint32       `json:"host_upstream_asn,omitempty" msgpack:"host_upstream_asn,omitempty" redis:"host_upstream_asn" redisearch:"numeric"`          // upstream route
	HostHostingUUID uint32       `json:"host_hosting_uuid,omitempty" msgpack:"host_hosting_uuid,omitempty" redis:"host_hosting_uuid" redisearch:"text"`             // (?) replace with member/memberOf
	HostURL         *mod_net.URL `json:"host_url,omitempty"          msgpack:"host_url,omitempty"          redis:"host_url"          redisearch:"text,sortable"`    //
	HostListen      *netip.Addr  `json:"host_listen,omitempty"       msgpack:"host_listen,omitempty"       redis:"host_listen"       redisearch:"text,sortable"`    //

	// specific data (space-separated KV DB stored as labeledURI)
	LabeledURI attrLabeledURIs `json:"labeledURI,omitempty" ldap:"labeledURI" msgpack:"labeledURI,omitempty" redis:"labeledURI" redisearch:"tag"` //
}

type AttrEntryType int
type AttrEntryStatus int
type entryFieldName string
