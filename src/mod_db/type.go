package mod_db

import (
	"net/netip"
	"net/url"

	"github.com/gomodule/redigo/redis"

	"rmm23/src/mod_ldap"
)

type Conn struct {
	conn redis.Conn
}

type EntryType int

// Entry is the struct that represents an LDAP-compatible entry.
type Entry struct {
	// element specific meta data
	Type EntryType `json:"type,omitempty" msgpack:"type,omitempty" redis:"type" redisearch:"text"` // entry's type `(domain|group|user|host)`

	// element meta data
	UUID            mod_ldap.AttrUUID          `json:"entryUUID,omitempty"       ldap:"entryUUID"       msgpack:"entryUUID,omitempty"       redis:"uuid"            redisearch:"text,sortable"` // must be unique
	DN              mod_ldap.AttrDN            `json:"dn,omitempty"              ldap:"dn"              msgpack:"dn,omitempty"              redis:"dn"              redisearch:"text,sortable"` // must be unique
	ObjectClass     mod_ldap.AttrObjectClasses `json:"objectClass,omitempty"     ldap:"objectClass"     msgpack:"objectClass,omitempty"     redis:"objectClass"     redisearch:"tag"`           // entry type
	CreatorsName    mod_ldap.AttrDN            `json:"creatorsName,omitempty"    ldap:"creatorsName"    msgpack:"creatorsName,omitempty"    redis:"creatorsName"    redisearch:"text"`          //
	CreateTimestamp mod_ldap.AttrTimestamp     `json:"createTimestamp,omitempty" ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"`          //
	ModifiersName   mod_ldap.AttrDN            `json:"modifiersName,omitempty"   ldap:"modifiersName"   msgpack:"modifiersName,omitempty"   redis:"modifiersName"   redisearch:"text"`          //
	ModifyTimestamp mod_ldap.AttrTimestamp     `json:"modifyTimestamp,omitempty" ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"`          //

	// element data
	CN                   mod_ldap.AttrString                `json:"cn,omitempty"                   ldap:"cn"                   msgpack:"cn,omitempty"                   redis:"cn"                   redisearch:"text"`             // RDN in group's context
	DC                   mod_ldap.AttrString                `json:"dc,omitempty"                   ldap:"dc"                   msgpack:"dc,omitempty"                   redis:"dc"                   redisearch:"text,sortable"`    //
	Description          mod_ldap.AttrString                `json:"description,omitempty"          ldap:"description"          msgpack:"description,omitempty"          redis:"description"          redisearch:"text"`             //
	DestinationIndicator mod_ldap.AttrDestinationIndicators `json:"destinationIndicator,omitempty" ldap:"destinationIndicator" msgpack:"destinationIndicator,omitempty" redis:"destinationIndicator" redisearch:"text"`             //
	DisplayName          mod_ldap.AttrString                `json:"displayName,omitempty"          ldap:"displayName"          msgpack:"displayName,omitempty"          redis:"displayName"          redisearch:"text,sortable"`    //
	GIDNumber            mod_ldap.AttrIDNumber              `json:"gidNumber,omitempty"            ldap:"gidNumber"            msgpack:"gidNumber,omitempty"            redis:"gidNumber"            redisearch:"numeric"`          // Primary GIDNumber in user's context (ignore it) and GIDNumber in group's context.
	HomeDirectory        mod_ldap.AttrString                `json:"homeDirectory,omitempty"        ldap:"homeDirectory"        msgpack:"homeDirectory,omitempty"        redis:"homeDirectory"        redisearch:"text"`             //
	IPHostNumber         mod_ldap.AttrIPHostNumbers         `json:"ipHostNumber,omitempty"         ldap:"ipHostNumber"         msgpack:"ipHostNumber,omitempty"         redis:"ipHostNumber"         redisearch:"text,sortable"`    //
	Mail                 mod_ldap.AttrMails                 `json:"mail,omitempty"                 ldap:"mail"                 msgpack:"mail,omitempty"                 redis:"mail"                 redisearch:"text"`             //
	Member               mod_ldap.AttrDNs                   `json:"member,omitempty"               ldap:"member"               msgpack:"member,omitempty"               redis:"member"               redisearch:"tag,sortable"`     //
	O                    mod_ldap.AttrString                `json:"o,omitempty"                    ldap:"o"                    msgpack:"o,omitempty"                    redis:"o"                    redisearch:"text"`             //
	OU                   mod_ldap.AttrString                `json:"ou,omitempty"                   ldap:"ou"                   msgpack:"ou,omitempty"                   redis:"ou"                   redisearch:"text"`             //
	Owner                mod_ldap.AttrDNs                   `json:"owner,omitempty"                ldap:"owner"                msgpack:"owner,omitempty"                redis:"owner"                redisearch:"tag"`              //
	SN                   mod_ldap.AttrString                `json:"sn,omitempty"                   ldap:"sn"                   msgpack:"sn,omitempty"                   redis:"sn"                   redisearch:"text"`             //
	SSHPublicKey         mod_ldap.AttrSSHPublicKeys         `json:"sshPublicKey,omitempty"         ldap:"sshPublicKey"         msgpack:"sshPublicKey,omitempty"         redis:"sshPublicKey"         redisearch:"tag"`              //
	TelephoneNumber      mod_ldap.AttrStrings               `json:"telephoneNumber,omitempty"      ldap:"telephoneNumber"      msgpack:"telephoneNumber,omitempty"      redis:"telephoneNumber"      redisearch:"text"`             //
	TelexNumber          mod_ldap.AttrStrings               `json:"telexNumber,omitempty"          ldap:"telexNumber"          msgpack:"telexNumber,omitempty"          redis:"telexNumber"          redisearch:"text"`             //
	UID                  mod_ldap.AttrID                    `json:"uid,omitempty"                  ldap:"uid"                  msgpack:"uid,omitempty"                  redis:"uid"                  redisearch:"text,sortable"`    // RDN in user's context
	UIDNumber            mod_ldap.AttrIDNumber              `json:"uidNumber,omitempty"            ldap:"uidNumber"            msgpack:"uidNumber,omitempty"            redis:"uidNumber"            redisearch:"numeric,sortable"` //
	UserPKCS12           mod_ldap.AttrUserPKCS12s           `json:"userPKCS12,omitempty"           ldap:"userPKCS12"           msgpack:"userPKCS12,omitempty"           redis:"userPKCS12"           redisearch:"tag"`              //
	UserPassword         mod_ldap.AttrUserPassword          `json:"userPassword,omitempty"         ldap:"userPassword"         msgpack:"userPassword,omitempty"         redis:"userPassword"         redisearch:"text"`             //
	// MemberOf             mod_ldap.AttrDNs                   `ldap:"memberOf" msgpack:"memberOf,omitempty" json:"memberOf,omitempty" redis:"memberOf" redisearch:"tag"`                                                  // ignore it, don't cache, calculate on the fly or avoid

	// specific data
	AAA string `json:"host_aaa,omitempty" msgpack:"host_aaa,omitempty" redis:"host_aaa" redisearch:"text"` // entry's AAA (?) `(UserPKCS12|UserPassword|SSHPublicKey|etc)`
	ACL string `json:"host_acl,omitempty" msgpack:"host_acl,omitempty" redis:"host_acl" redisearch:"text"` // entry's ACL

	// host specific data
	HostType        string     `json:"host_type,omitempty"         msgpack:"host_type,omitempty"         redis:"host_type"         redisearch:"text"`             // host type `(provider|interim|openvpn|ciscovpn)`
	HostASN         uint32     `json:"host_asn,omitempty"          msgpack:"host_asn,omitempty"          redis:"host_asn"          redisearch:"numeric,sortable"` //
	HostUpstreamASN uint32     `json:"host_upstream_asn,omitempty" msgpack:"host_upstream_asn,omitempty" redis:"host_upstream_asn" redisearch:"numeric"`          // upstream route
	HostHostingUUID uint32     `json:"host_hosting_uuid,omitempty" msgpack:"host_hosting_uuid,omitempty" redis:"host_hosting_uuid" redisearch:"text"`             // (?) replace with member/memberOf
	HostURL         url.URL    `json:"host_url,omitempty"          msgpack:"host_url,omitempty"          redis:"host_url"          redisearch:"text,sortable"`    //
	HostListen      netip.Addr `json:"host_listen,omitempty"       msgpack:"host_listen,omitempty"       redis:"host_listen"       redisearch:"text,sortable"`    //

	// specific data (space-separated KV DB stored as labeledURIs)
	LabeledURI mod_ldap.AttrLabeledURIs `json:"labeledURI,omitempty" ldap:"labeledURI" msgpack:"labeledURI,omitempty" redis:"labeledURI" redisearch:"tag"` //
}

type Domain struct{ Entry }
type Group struct{ Entry }
type User struct{ Entry }
type Host struct{ Entry }
