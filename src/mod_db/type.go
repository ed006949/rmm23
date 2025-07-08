package mod_db

import (
	"net/netip"
	"net/url"

	"rmm23/src/mod_ldap"
)

type ElementEntry struct {
	// element meta
	UUID            mod_ldap.AttrUUID          `ldap:"entryUUID" msgpack:"entryUUID,omitempty" redis:"uuid" redisearch:"text,sortable"`               // must be unique
	DN              mod_ldap.AttrDN            `ldap:"dn" msgpack:"dn,omitempty" redis:"dn" redisearch:"text,sortable"`                               // must be unique
	ObjectClass     mod_ldap.AttrObjectClasses `ldap:"objectClass" msgpack:"objectClass,omitempty" redis:"objectClass" redisearch:"text"`             // entry type
	CreatorsName    mod_ldap.AttrDN            `ldap:"creatorsName" msgpack:"creatorsName,omitempty" redis:"creatorsName" redisearch:"text"`          //
	CreateTimestamp mod_ldap.AttrTimestamp     `ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"` //
	ModifiersName   mod_ldap.AttrDN            `ldap:"modifiersName" msgpack:"modifiersName,omitempty" redis:"modifiersName" redisearch:"text"`       //
	ModifyTimestamp mod_ldap.AttrTimestamp     `ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"` //

	// element data
	CN                   mod_ldap.AttrString                `ldap:"cn" msgpack:"cn,omitempty" redis:"cn" redisearch:"text"`                                                       // RDN in group's context
	DC                   mod_ldap.AttrString                `ldap:"dc" msgpack:"dc,omitempty" redis:"dc" redisearch:"text,sortable"`                                              //
	Description          mod_ldap.AttrString                `ldap:"description" msgpack:"description,omitempty" redis:"description" redisearch:"text"`                            //
	DestinationIndicator mod_ldap.AttrDestinationIndicators `ldap:"destinationIndicator" msgpack:"destinationIndicator,omitempty" redis:"destinationIndicator" redisearch:"text"` //
	DisplayName          mod_ldap.AttrString                `ldap:"displayName" msgpack:"displayName,omitempty" redis:"displayName" redisearch:"text,sortable"`                   //
	GIDNumber            mod_ldap.AttrIDNumber              `ldap:"gidNumber" msgpack:"gidNumber,omitempty" redis:"gidNumber" redisearch:"numeric"`                               // Primary GIDNumber in user's context and GIDNumber in group's context
	HomeDirectory        mod_ldap.AttrString                `ldap:"homeDirectory" msgpack:"homeDirectory,omitempty" redis:"homeDirectory" redisearch:"text"`                      //
	IPHostNumber         mod_ldap.AttrIPHostNumbers         `ldap:"ipHostNumber" msgpack:"ipHostNumber,omitempty" redis:"ipHostNumber" redisearch:"text,sortable"`                //
	Mail                 mod_ldap.AttrMails                 `ldap:"mail" msgpack:"mail,omitempty" redis:"mail" redisearch:"text"`                                                 //
	Member               mod_ldap.AttrDNs                   `ldap:"member" msgpack:"member,omitempty" redis:"member" redisearch:"text,sortable"`                                  //
	MemberOf             mod_ldap.AttrDNs                   `ldap:"memberOf" msgpack:"memberOf,omitempty" redis:"memberOf" redisearch:"text"`                                     //
	O                    mod_ldap.AttrString                `ldap:"o" msgpack:"o,omitempty" redis:"o" redisearch:"text"`                                                          //
	OU                   mod_ldap.AttrString                `ldap:"ou" msgpack:"ou,omitempty" redis:"ou" redisearch:"text"`                                                       //
	Owner                mod_ldap.AttrDNs                   `ldap:"owner" msgpack:"owner,omitempty" redis:"owner" redisearch:"text"`                                              //
	SN                   mod_ldap.AttrString                `ldap:"sn" msgpack:"sn,omitempty" redis:"sn" redisearch:"text"`                                                       //
	SSHPublicKey         mod_ldap.AttrSSHPublicKeys         `ldap:"sshPublicKey" msgpack:"sshPublicKey,omitempty" redis:"sshPublicKey" redisearch:"text"`                         //
	TelephoneNumber      mod_ldap.AttrStrings               `ldap:"telephoneNumber" msgpack:"telephoneNumber,omitempty" redis:"telephoneNumber" redisearch:"text"`                //
	TelexNumber          mod_ldap.AttrStrings               `ldap:"telexNumber" msgpack:"telexNumber,omitempty" redis:"telexNumber" redisearch:"text"`                            //
	UID                  mod_ldap.AttrID                    `ldap:"uid" msgpack:"uid,omitempty" redis:"uid" redisearch:"text,sortable"`                                           // RDN in user's context
	UIDNumber            mod_ldap.AttrIDNumber              `ldap:"uidNumber" msgpack:"uidNumber,omitempty" redis:"uidNumber" redisearch:"numeric,sortable"`                      //
	UserPKCS12           mod_ldap.AttrUserPKCS12s           `ldap:"userPKCS12" msgpack:"userPKCS12,omitempty" redis:"userPKCS12" redisearch:"text"`                               //
	UserPassword         mod_ldap.AttrUserPassword          `ldap:"userPassword" msgpack:"userPassword,omitempty" redis:"userPassword" redisearch:"text"`                         //

	// host specific data
	Type        string     `xml:"type,attr,omitempty" msgpack:"type,omitempty" redis:"host_type" redisearch:"text"`                            // host type `(provider|interim|openvpn|ciscovpn)`
	ASN         uint32     `xml:"asn,attr,omitempty" msgpack:"asn,omitempty" redis:"host_asn" redisearch:"numeric,sortable"`                   //
	UpstreamASN uint32     `xml:"upstream_asn,attr,omitempty" msgpack:"upstream_asn,omitempty" redis:"host_upstream_asn" redisearch:"numeric"` //
	HostUUID    uint32     `xml:"host_uuid,attr,omitempty" msgpack:"host_uuid,omitempty" redis:"host_host_uuid" redisearch:"numeric"`          // (?) replace with member/memberOf
	URL         url.URL    `xml:"url,attr,omitempty"`                                                                                          //
	Listen      netip.Addr `xml:"listen,attr,omitempty"`                                                                                       //
	ACL         string     `xml:"acl,attr,omitempty"`                                                                                          //
	AAA         string     `xml:"aaa,attr,omitempty"`                                                                                          //
}

type ElementDomain struct {
	UUID            mod_ldap.AttrUUID          `ldap:"entryUUID" msgpack:"entryUUID,omitempty" redis:"uuid" redisearch:"text,sortable"`
	DN              mod_ldap.AttrDN            `ldap:"dn" msgpack:"dn,omitempty" redis:"dn" redisearch:"text,sortable"`
	ObjectClass     mod_ldap.AttrObjectClasses `ldap:"objectClass" msgpack:"objectClass,omitempty" redis:"objectClass" redisearch:"text"`
	CreatorsName    mod_ldap.AttrDN            `ldap:"creatorsName" msgpack:"creatorsName,omitempty" redis:"creatorsName" redisearch:"text"`
	CreateTimestamp mod_ldap.AttrTimestamp     `ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"`
	ModifiersName   mod_ldap.AttrDN            `ldap:"modifiersName" msgpack:"modifiersName,omitempty" redis:"modifiersName" redisearch:"text"`
	ModifyTimestamp mod_ldap.AttrTimestamp     `ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"`

	DC mod_ldap.AttrString `ldap:"dc" msgpack:"dc,omitempty" redis:"dc" redisearch:"text,sortable"`
	O  mod_ldap.AttrString `ldap:"o" msgpack:"o,omitempty" redis:"o" redisearch:"text,sortable"`
}
type ElementGroup struct {
	UUID            mod_ldap.AttrUUID          `ldap:"entryUUID" msgpack:"entryUUID,omitempty" redis:"uuid" redisearch:"text,sortable"`
	DN              mod_ldap.AttrDN            `ldap:"dn" msgpack:"dn,omitempty" redis:"dn" redisearch:"text,sortable"`
	ObjectClass     mod_ldap.AttrObjectClasses `ldap:"objectClass" msgpack:"objectClass,omitempty" redis:"objectClass" redisearch:"text"`
	CreatorsName    mod_ldap.AttrDN            `ldap:"creatorsName" msgpack:"creatorsName,omitempty" redis:"creatorsName" redisearch:"text"`
	CreateTimestamp mod_ldap.AttrTimestamp     `ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"`
	ModifiersName   mod_ldap.AttrDN            `ldap:"modifiersName" msgpack:"modifiersName,omitempty" redis:"modifiersName" redisearch:"text"`
	ModifyTimestamp mod_ldap.AttrTimestamp     `ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"`

	CN        mod_ldap.AttrString   `ldap:"cn" msgpack:"cn,omitempty" redis:"cn" redisearch:"text"`
	GIDNumber mod_ldap.AttrIDNumber `ldap:"gidNumber" msgpack:"gidNumber,omitempty" redis:"gidNumber" redisearch:"numeric,sortable"`
	Member    mod_ldap.AttrDNs      `ldap:"member" msgpack:"member,omitempty" redis:"member" redisearch:"text,sortable"`
	Owner     mod_ldap.AttrDNs      `ldap:"owner" msgpack:"owner,omitempty" redis:"owner" redisearch:"text"`

	LabeledURI mod_ldap.AttrLabeledURIs `ldap:"labeledURI"`
}
type ElementUser struct {
	UUID            mod_ldap.AttrUUID          `ldap:"entryUUID" msgpack:"entryUUID,omitempty" redis:"uuid" redisearch:"text,sortable"`
	DN              mod_ldap.AttrDN            `ldap:"dn" msgpack:"dn,omitempty" redis:"dn" redisearch:"text,sortable"`
	ObjectClass     mod_ldap.AttrObjectClasses `ldap:"objectClass" msgpack:"objectClass,omitempty" redis:"objectClass" redisearch:"text"`
	CreatorsName    mod_ldap.AttrDN            `ldap:"creatorsName" msgpack:"creatorsName,omitempty" redis:"creatorsName" redisearch:"text"`
	CreateTimestamp mod_ldap.AttrTimestamp     `ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"`
	ModifiersName   mod_ldap.AttrDN            `ldap:"modifiersName" msgpack:"modifiersName,omitempty" redis:"modifiersName" redisearch:"text"`
	ModifyTimestamp mod_ldap.AttrTimestamp     `ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"`

	CN                   mod_ldap.AttrString                `ldap:"cn" msgpack:"cn,omitempty" redis:"cn" redisearch:"text"`
	Description          mod_ldap.AttrString                `ldap:"description" msgpack:"description,omitempty" redis:"description" redisearch:"text"`
	DestinationIndicator mod_ldap.AttrDestinationIndicators `ldap:"destinationIndicator" msgpack:"destinationIndicator,omitempty" redis:"destinationIndicator" redisearch:"text"`
	DisplayName          mod_ldap.AttrString                `ldap:"displayName" msgpack:"displayName,omitempty" redis:"displayName" redisearch:"text,sortable"`
	GIDNumber            mod_ldap.AttrIDNumber              `ldap:"gidNumber" msgpack:"gidNumber,omitempty" redis:"gidNumber" redisearch:"numeric"`
	HomeDirectory        mod_ldap.AttrString                `ldap:"homeDirectory" msgpack:"homeDirectory,omitempty" redis:"homeDirectory" redisearch:"text"`
	IPHostNumber         mod_ldap.AttrIPHostNumbers         `ldap:"ipHostNumber" msgpack:"ipHostNumber,omitempty" redis:"ipHostNumber" redisearch:"text,sortable"`
	Mail                 mod_ldap.AttrMails                 `ldap:"mail" msgpack:"mail,omitempty" redis:"mail" redisearch:"text"`
	MemberOf             mod_ldap.AttrDNs                   `ldap:"memberOf" msgpack:"memberOf,omitempty"`
	O                    mod_ldap.AttrString                `ldap:"o" msgpack:"o,omitempty" redis:"o" redisearch:"text"`
	OU                   mod_ldap.AttrString                `ldap:"ou" msgpack:"ou,omitempty" redis:"ou" redisearch:"text"`
	SN                   mod_ldap.AttrString                `ldap:"sn" msgpack:"sn,omitempty" redis:"sn" redisearch:"text"`
	SSHPublicKey         mod_ldap.AttrSSHPublicKeys         `ldap:"sshPublicKey" msgpack:"sshPublicKey,omitempty" redis:"sshPublicKey" redisearch:"text"`
	TelephoneNumber      mod_ldap.AttrStrings               `ldap:"telephoneNumber" msgpack:"telephoneNumber,omitempty" redis:"telephoneNumber" redisearch:"text"`
	TelexNumber          mod_ldap.AttrStrings               `ldap:"telexNumber" msgpack:"telexNumber,omitempty" redis:"telexNumber" redisearch:"text"`
	UID                  mod_ldap.AttrID                    `ldap:"uid" msgpack:"uid,omitempty" redis:"uid" redisearch:"text,sortable"`
	UIDNumber            mod_ldap.AttrIDNumber              `ldap:"uidNumber" msgpack:"uidNumber,omitempty" redis:"uidNumber" redisearch:"numeric,sortable"`
	UserPKCS12           mod_ldap.AttrUserPKCS12s           `ldap:"userPKCS12" msgpack:"userPKCS12,omitempty" redis:"userPKCS12" redisearch:"text"`
	UserPassword         mod_ldap.AttrUserPassword          `ldap:"userPassword" msgpack:"userPassword,omitempty" redis:"userPassword" redisearch:"text"`

	LabeledURI mod_ldap.AttrLabeledURIs `ldap:"labeledURI"`
}
type ElementHost struct {
	UUID            mod_ldap.AttrUUID          `ldap:"entryUUID" msgpack:"entryUUID,omitempty" redis:"uuid" redisearch:"text,sortable"`
	DN              mod_ldap.AttrDN            `ldap:"dn" msgpack:"dn,omitempty" redis:"dn" redisearch:"text,sortable"`
	ObjectClass     mod_ldap.AttrObjectClasses `ldap:"objectClass" msgpack:"objectClass,omitempty" redis:"objectClass" redisearch:"text"`
	CreatorsName    mod_ldap.AttrDN            `ldap:"creatorsName" msgpack:"creatorsName,omitempty" redis:"creatorsName" redisearch:"text"`
	CreateTimestamp mod_ldap.AttrTimestamp     `ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"`
	ModifiersName   mod_ldap.AttrDN            `ldap:"modifiersName" msgpack:"modifiersName,omitempty" redis:"modifiersName" redisearch:"text"`
	ModifyTimestamp mod_ldap.AttrTimestamp     `ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"`

	CN            mod_ldap.AttrString      `ldap:"cn" msgpack:"cn,omitempty" redis:"cn" redisearch:"text"`
	GIDNumber     mod_ldap.AttrIDNumber    `ldap:"gidNumber" msgpack:"gidNumber,omitempty" redis:"gidNumber" redisearch:"numeric"`
	HomeDirectory mod_ldap.AttrString      `ldap:"homeDirectory" msgpack:"homeDirectory,omitempty" redis:"homeDirectory" redisearch:"text"`
	SN            mod_ldap.AttrString      `ldap:"sn" msgpack:"sn,omitempty" redis:"sn" redisearch:"text"`
	UID           mod_ldap.AttrID          `ldap:"uid" msgpack:"uid,omitempty" redis:"uid" redisearch:"text,sortable"`
	UIDNumber     mod_ldap.AttrIDNumber    `ldap:"uidNumber" msgpack:"uidNumber,omitempty" redis:"uidNumber" redisearch:"numeric,sortable"`
	UserPKCS12    mod_ldap.AttrUserPKCS12s `ldap:"userPKCS12" msgpack:"userPKCS12,omitempty" redis:"userPKCS12" redisearch:"text"`

	LabeledURI mod_ldap.AttrLabeledURIs `ldap:"labeledURI"`

	// Type `(provider|interim|openvpn|ciscovpn)`
	Type        string     `xml:"type,attr,omitempty" msgpack:"type,omitempty" redis:"host_type" redisearch:"text"`
	ASN         uint32     `xml:"asn,attr,omitempty" msgpack:"asn,omitempty" redis:"host_asn" redisearch:"numeric,sortable"`
	UpstreamASN uint32     `xml:"upstream_asn,attr,omitempty" msgpack:"upstream_asn,omitempty" redis:"host_upstream_asn" redisearch:"numeric"`
	HostUUID    uint32     `xml:"host_uuid,attr,omitempty" msgpack:"host_uuid,omitempty" redis:"host_host_uuid" redisearch:"numeric"`
	URL         url.URL    `xml:"url,attr,omitempty" msgpack:"url,omitempty"`
	Listen      netip.Addr `xml:"listen,attr,omitempty" msgpack:"listen,omitempty"`
	ACL         string     `xml:"acl,attr,omitempty" msgpack:"acl,omitempty"`
	AAA         string     `xml:"aaa,attr,omitempty" msgpack:"aaa,omitempty"`

	// (?)
	Member   mod_ldap.AttrDNs `ldap:"member" msgpack:"member,omitempty" redis:"member" redisearch:"text,sortable"`
	Owner    mod_ldap.AttrDNs `ldap:"owner" msgpack:"owner,omitempty" redis:"owner" redisearch:"text"`
	MemberOf mod_ldap.AttrDNs `ldap:"memberOf" msgpack:"memberOf,omitempty" redis:"memberOf" redisearch:"text"`
}
