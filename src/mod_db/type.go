package mod_db

import (
	"net/netip"
	"net/url"

	"rmm23/src/mod_ldap"
)

type EntryType int

// Entry is the struct that represents an LDAP-compatible entry.
type Entry struct {
	// element meta data
	UUID            mod_ldap.AttrUUID          `ldap:"entryUUID" msgpack:"entryUUID,omitempty" json:"entryUUID,omitempty" redis:"uuid" redisearch:"text,sortable"`                     // must be unique
	DN              mod_ldap.AttrDN            `ldap:"dn" msgpack:"dn,omitempty" json:"dn,omitempty" redis:"dn" redisearch:"text,sortable"`                                            // must be unique
	ObjectClass     mod_ldap.AttrObjectClasses `ldap:"objectClass" msgpack:"objectClass,omitempty" json:"objectClass,omitempty" redis:"objectClass" redisearch:"text"`                 // entry type
	CreatorsName    mod_ldap.AttrDN            `ldap:"creatorsName" msgpack:"creatorsName,omitempty" json:"creatorsName,omitempty" redis:"creatorsName" redisearch:"text"`             //
	CreateTimestamp mod_ldap.AttrTimestamp     `ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" json:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"` //
	ModifiersName   mod_ldap.AttrDN            `ldap:"modifiersName" msgpack:"modifiersName,omitempty" json:"modifiersName,omitempty" redis:"modifiersName" redisearch:"text"`         //
	ModifyTimestamp mod_ldap.AttrTimestamp     `ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" json:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"` //

	// element data
	CN                   mod_ldap.AttrString                `ldap:"cn" msgpack:"cn,omitempty" json:"cn,omitempty" redis:"cn" redisearch:"text"`                                                                         // RDN in group's context
	DC                   mod_ldap.AttrString                `ldap:"dc" msgpack:"dc,omitempty" json:"dc,omitempty" redis:"dc" redisearch:"text,sortable"`                                                                //
	Description          mod_ldap.AttrString                `ldap:"description" msgpack:"description,omitempty" json:"description,omitempty" redis:"description" redisearch:"text"`                                     //
	DestinationIndicator mod_ldap.AttrDestinationIndicators `ldap:"destinationIndicator" msgpack:"destinationIndicator,omitempty" json:"destinationIndicator,omitempty" redis:"destinationIndicator" redisearch:"text"` //
	DisplayName          mod_ldap.AttrString                `ldap:"displayName" msgpack:"displayName,omitempty" json:"displayName,omitempty" redis:"displayName" redisearch:"text,sortable"`                            //
	GIDNumber            mod_ldap.AttrIDNumber              `ldap:"gidNumber" msgpack:"gidNumber,omitempty" json:"gidNumber,omitempty" redis:"gidNumber" redisearch:"numeric"`                                          // Primary GIDNumber in user's context (ignore it) and GIDNumber in group's context.
	HomeDirectory        mod_ldap.AttrString                `ldap:"homeDirectory" msgpack:"homeDirectory,omitempty" json:"homeDirectory,omitempty" redis:"homeDirectory" redisearch:"text"`                             //
	IPHostNumber         mod_ldap.AttrIPHostNumbers         `ldap:"ipHostNumber" msgpack:"ipHostNumber,omitempty" json:"ipHostNumber,omitempty" redis:"ipHostNumber" redisearch:"text,sortable"`                        //
	Mail                 mod_ldap.AttrMails                 `ldap:"mail" msgpack:"mail,omitempty" json:"mail,omitempty" redis:"mail" redisearch:"text"`                                                                 //
	Member               mod_ldap.AttrDNs                   `ldap:"member" msgpack:"member,omitempty" json:"member,omitempty" redis:"member" redisearch:"text,sortable"`                                                //
	MemberOf             mod_ldap.AttrDNs                   `ldap:"memberOf" msgpack:"memberOf,omitempty" json:"memberOf,omitempty" redis:"memberOf" redisearch:"text"`                                                 // ignore it, don't cache, calculate on the fly or avoid
	O                    mod_ldap.AttrString                `ldap:"o" msgpack:"o,omitempty" json:"o,omitempty" redis:"o" redisearch:"text"`                                                                             //
	OU                   mod_ldap.AttrString                `ldap:"ou" msgpack:"ou,omitempty" json:"ou,omitempty" redis:"ou" redisearch:"text"`                                                                         //
	Owner                mod_ldap.AttrDNs                   `ldap:"owner" msgpack:"owner,omitempty" json:"owner,omitempty" redis:"owner" redisearch:"text"`                                                             //
	SN                   mod_ldap.AttrString                `ldap:"sn" msgpack:"sn,omitempty" json:"sn,omitempty" redis:"sn" redisearch:"text"`                                                                         //
	SSHPublicKey         mod_ldap.AttrSSHPublicKeys         `ldap:"sshPublicKey" msgpack:"sshPublicKey,omitempty" json:"sshPublicKey,omitempty" redis:"sshPublicKey" redisearch:"text"`                                 //
	TelephoneNumber      mod_ldap.AttrStrings               `ldap:"telephoneNumber" msgpack:"telephoneNumber,omitempty" json:"telephoneNumber,omitempty" redis:"telephoneNumber" redisearch:"text"`                     //
	TelexNumber          mod_ldap.AttrStrings               `ldap:"telexNumber" msgpack:"telexNumber,omitempty" json:"telexNumber,omitempty" redis:"telexNumber" redisearch:"text"`                                     //
	UID                  mod_ldap.AttrID                    `ldap:"uid" msgpack:"uid,omitempty" json:"uid,omitempty" redis:"uid" redisearch:"text,sortable"`                                                            // RDN in user's context
	UIDNumber            mod_ldap.AttrIDNumber              `ldap:"uidNumber" msgpack:"uidNumber,omitempty" json:"uidNumber,omitempty" redis:"uidNumber" redisearch:"numeric,sortable"`                                 //
	UserPKCS12           mod_ldap.AttrUserPKCS12s           `ldap:"userPKCS12" msgpack:"userPKCS12,omitempty" json:"userPKCS12,omitempty" redis:"userPKCS12" redisearch:"text"`                                         //
	UserPassword         mod_ldap.AttrUserPassword          `ldap:"userPassword" msgpack:"userPassword,omitempty" json:"userPassword,omitempty" redis:"userPassword" redisearch:"text"`                                 //

	// specific data (space-separated KV DB stored as labelledURIs)
	Legacy mod_ldap.LabeledURILegacy `ldap:"labelledURI" msgpack:"-" json:"-" redis:"-" redisearch:"-"` //

	// specific data
	Type EntryType `msgpack:"type,omitempty" json:"type,omitempty" redis:"type" redisearch:"text"`             // entry's type `(domain|group|user|host)`
	AAA  string    `msgpack:"host_aaa,omitempty" json:"host_aaa,omitempty" redis:"host_aaa" redisearch:"text"` // entry's AAA (?) `(UserPKCS12|UserPassword|SSHPublicKey|etc)`
	ACL  string    `msgpack:"host_acl,omitempty" json:"host_acl,omitempty" redis:"host_acl" redisearch:"text"` // entry's ACL

	// host specific data
	HostType        string     `msgpack:"host_type,omitempty" json:"host_type,omitempty" redis:"host_type" redisearch:"text"`                            // host type `(provider|interim|openvpn|ciscovpn)`
	HostASN         uint32     `msgpack:"host_asn,omitempty" json:"host_asn,omitempty" redis:"host_asn" redisearch:"numeric,sortable"`                   //
	HostUpstreamASN uint32     `msgpack:"host_upstream_asn,omitempty" json:"host_upstream_asn,omitempty" redis:"host_upstream_asn" redisearch:"numeric"` // upstream route
	HostHostingUUID uint32     `msgpack:"host_hosting_uuid,omitempty" json:"host_hosting_uuid,omitempty" redis:"host_hosting_uuid" redisearch:"text"`    // (?) replace with member/memberOf
	HostURL         url.URL    `msgpack:"host_url,omitempty" json:"host_url,omitempty" redis:"host_url" redisearch:"text,sortable"`                      //
	HostListen      netip.Addr `msgpack:"host_listen,omitempty" json:"host_listen,omitempty" redis:"host_listen" redisearch:"text,sortable"`             //
}
