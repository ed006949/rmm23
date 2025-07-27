package mod_db

import (
	"net/netip"

	"rmm23/src/mod_net"
)

// entry is the struct that represents an LDAP-compatible entry.
//
// when updating @src/mod_db/entry_type.go don't forget to update @src/mod_db/entry_const.go.
type entry struct {
	// element specific meta data
	Type   attrEntryType   `json:"type,omitempty"   msgpack:"type,omitempty"   redis:"type"   redisearch:"numeric,sortable"` // entry's type `(domain|group|user|host)`
	Status attrEntryStatus `json:"status,omitempty" msgpack:"status,omitempty" redis:"status" redisearch:"numeric,sortable"` //
	BaseDN attrDN          `json:"baseDN,omitempty" msgpack:"baseDN,omitempty" redis:"baseDN" redisearch:"tag,sortable"`     //

	// element meta data
	UUID            attrUUID          `json:"uuid,omitempty"            ldap:"entryUUID"       msgpack:"uuid,omitempty"            redis:"uuid"            redisearch:"tag,sortable"` // must be unique
	DN              attrDN            `json:"dn,omitempty"              ldap:"dn"              msgpack:"dn,omitempty"              redis:"dn"              redisearch:"tag,sortable"` // must be unique
	ObjectClass     attrObjectClasses `json:"objectClass,omitempty"     ldap:"objectClass"     msgpack:"objectClass,omitempty"     redis:"objectClass"     redisearch:"tag"`          // entry type
	CreatorsName    attrDN            `json:"creatorsName,omitempty"    ldap:"creatorsName"    msgpack:"creatorsName,omitempty"    redis:"creatorsName"    redisearch:"tag"`          //
	CreateTimestamp attrTimestamp     `json:"createTimestamp,omitempty" ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"tag"`          //
	ModifiersName   attrDN            `json:"modifiersName,omitempty"   ldap:"modifiersName"   msgpack:"modifiersName,omitempty"   redis:"modifiersName"   redisearch:"tag"`          //
	ModifyTimestamp attrTimestamp     `json:"modifyTimestamp,omitempty" ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"tag"`          //

	// element data
	CN                   attrString                `json:"cn,omitempty"                   ldap:"cn"                   msgpack:"cn,omitempty"                   redis:"cn"                   redisearch:"tag"`              // RDN in group's context
	DC                   attrString                `json:"dc,omitempty"                   ldap:"dc"                   msgpack:"dc,omitempty"                   redis:"dc"                   redisearch:"tag,sortable"`     //
	Description          attrString                `json:"description,omitempty"          ldap:"description"          msgpack:"description,omitempty"          redis:"description"          redisearch:"tag"`              //
	DestinationIndicator attrDestinationIndicators `json:"destinationIndicator,omitempty" ldap:"destinationIndicator" msgpack:"destinationIndicator,omitempty" redis:"destinationIndicator" redisearch:"tag"`              //
	DisplayName          attrString                `json:"displayName,omitempty"          ldap:"displayName"          msgpack:"displayName,omitempty"          redis:"displayName"          redisearch:"tag,sortable"`     //
	GIDNumber            attrIDNumber              `json:"gidNumber,omitempty"            ldap:"gidNumber"            msgpack:"gidNumber,omitempty"            redis:"gidNumber"            redisearch:"numeric,sortable"` // Primary GIDNumber in user's context (ignore it) and GIDNumber in group's context.
	HomeDirectory        attrString                `json:"homeDirectory,omitempty"        ldap:"homeDirectory"        msgpack:"homeDirectory,omitempty"        redis:"homeDirectory"        redisearch:"tag"`              //
	IPHostNumber         attrIPHostNumbers         `json:"ipHostNumber,omitempty"         ldap:"ipHostNumber"         msgpack:"ipHostNumber,omitempty"         redis:"ipHostNumber"         redisearch:"tag,sortable"`     //
	Mail                 attrMails                 `json:"mail,omitempty"                 ldap:"mail"                 msgpack:"mail,omitempty"                 redis:"mail"                 redisearch:"tag"`              //
	Member               attrDNs                   `json:"member,omitempty"               ldap:"member"               msgpack:"member,omitempty"               redis:"member"               redisearch:"tag,sortable"`     //
	O                    attrString                `json:"o,omitempty"                    ldap:"o"                    msgpack:"o,omitempty"                    redis:"o"                    redisearch:"tag"`              //
	OU                   attrString                `json:"ou,omitempty"                   ldap:"ou"                   msgpack:"ou,omitempty"                   redis:"ou"                   redisearch:"tag"`              //
	Owner                attrDNs                   `json:"owner,omitempty"                ldap:"owner"                msgpack:"owner,omitempty"                redis:"owner"                redisearch:"tag"`              //
	SN                   attrString                `json:"sn,omitempty"                   ldap:"sn"                   msgpack:"sn,omitempty"                   redis:"sn"                   redisearch:"tag"`              //
	SSHPublicKey         attrSSHPublicKeys         `json:"sshPublicKey,omitempty"         ldap:"sshPublicKey"         msgpack:"sshPublicKey,omitempty"         redis:"sshPublicKey"         redisearch:"tag"`              //
	TelephoneNumber      attrStrings               `json:"telephoneNumber,omitempty"      ldap:"telephoneNumber"      msgpack:"telephoneNumber,omitempty"      redis:"telephoneNumber"      redisearch:"tag"`              //
	TelexNumber          attrStrings               `json:"telexNumber,omitempty"          ldap:"telexNumber"          msgpack:"telexNumber,omitempty"          redis:"telexNumber"          redisearch:"tag"`              //
	UID                  attrID                    `json:"uid,omitempty"                  ldap:"uid"                  msgpack:"uid,omitempty"                  redis:"uid"                  redisearch:"tag,sortable"`     // RDN in user's context
	UIDNumber            attrIDNumber              `json:"uidNumber,omitempty"            ldap:"uidNumber"            msgpack:"uidNumber,omitempty"            redis:"uidNumber"            redisearch:"numeric,sortable"` //
	UserPKCS12           attrUserPKCS12s           `json:"userPKCS12,omitempty"           ldap:"userPKCS12"           msgpack:"userPKCS12,omitempty"           redis:"userPKCS12"           redisearch:"tag"`              //
	UserPassword         attrUserPassword          `json:"userPassword,omitempty"         ldap:"userPassword"         msgpack:"userPassword,omitempty"         redis:"userPassword"         redisearch:"tag"`              //
	// MemberOf             attrDNs                   `json:"memberOf,omitempty"             ldap:"memberOf"             msgpack:"memberOf,omitempty"             redis:"memberOf"             redisearch:"tag"`              // ignore it, don't cache, calculate on the fly or avoid

	// specific data
	AAA string `json:"host_aaa,omitempty" msgpack:"host_aaa,omitempty" redis:"host_aaa" redisearch:"tag"` // entry's AAA (?) `(UserPKCS12|UserPassword|SSHPublicKey|etc)`
	ACL string `json:"host_acl,omitempty" msgpack:"host_acl,omitempty" redis:"host_acl" redisearch:"tag"` // entry's ACL

	// host specific data
	HostType        string       `json:"host_type,omitempty"         msgpack:"host_type,omitempty"         redis:"host_type"         redisearch:"tag"`          // host type `(provider|interim|openvpn|ciscovpn)`
	HostASN         uint32       `json:"host_asn,omitempty"          msgpack:"host_asn,omitempty"          redis:"host_asn"          redisearch:"tag,sortable"` //
	HostUpstreamASN uint32       `json:"host_upstream_asn,omitempty" msgpack:"host_upstream_asn,omitempty" redis:"host_upstream_asn" redisearch:"tag"`          // upstream route
	HostHostingUUID uint32       `json:"host_hosting_uuid,omitempty" msgpack:"host_hosting_uuid,omitempty" redis:"host_hosting_uuid" redisearch:"tag"`          // (?) replace with member/memberOf
	HostURL         *mod_net.URL `json:"host_url,omitempty"          msgpack:"host_url,omitempty"          redis:"host_url"          redisearch:"tag,sortable"` //
	HostListen      *netip.Addr  `json:"host_listen,omitempty"       msgpack:"host_listen,omitempty"       redis:"host_listen"       redisearch:"tag,sortable"` //

	// specific data (space-separated KV DB stored as labeledURI)
	LabeledURI attrLabeledURIs `json:"labeledURI,omitempty" ldap:"labeledURI" msgpack:"labeledURI,omitempty" redis:"labeledURI" redisearch:"tag"` //
}
