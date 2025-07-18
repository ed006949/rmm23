package mod_ldap

import (
	"net/netip"
	"net/url"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/mod_crypto"
	"rmm23/src/mod_net"
	"rmm23/src/mod_ssh"
)

type entries struct {
	url   *mod_net.URL
	conn  *ldap.Conn
	entry map[AttrDN]*ldap.Entry
}

type Conf struct {
	URL      *mod_net.URL    `xml:"url,attr"`
	Settings []*ConfSettings `xml:"settings"`
	Domain   []*ConfDomain   `xml:"domain"`

	// schema map[string]*schema
	conn *ldap.Conn

	Table *ConfTable
}
type ConfSettings struct {
	Type   string `xml:"type,attr"`
	DN     AttrDN `xml:"dn,attr"`
	CN     string `xml:"cn,attr"`
	Filter string `xml:"filter,attr"`
}
type ConfDomain struct {
	DN AttrDN `xml:"dn,attr"`

	Domain *Element
	Users  Elements
	Groups Elements
	Hosts  Elements

	searchResults map[string]*ldap.SearchResult
}
type ConfTable struct {
	Domain       Elements
	Users        Elements
	Groups       Elements
	Hosts        Elements
	IPHostNumber map[netip.Prefix]struct{}
	ID           map[AttrIDNumber]struct{}
}

type Elements map[AttrDN]*Element

type Element struct {
	UUID            AttrUUID          `ldap:"entryUUID" msgpack:"entryUUID,omitempty" json:"entryUUID,omitempty" redis:"uuid" redisearch:"text,sortable"`                     // must be unique
	DN              AttrDN            `ldap:"dn" msgpack:"dn,omitempty" json:"dn,omitempty" redis:"dn" redisearch:"text,sortable"`                                            // must be unique
	ObjectClass     AttrObjectClasses `ldap:"objectClass" msgpack:"objectClass,omitempty" json:"objectClass,omitempty" redis:"objectClass" redisearch:"tag"`                  // entry type
	CreatorsName    AttrDN            `ldap:"creatorsName" msgpack:"creatorsName,omitempty" json:"creatorsName,omitempty" redis:"creatorsName" redisearch:"text"`             //
	CreateTimestamp AttrTimestamp     `ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" json:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"` //
	ModifiersName   AttrDN            `ldap:"modifiersName" msgpack:"modifiersName,omitempty" json:"modifiersName,omitempty" redis:"modifiersName" redisearch:"text"`         //
	ModifyTimestamp AttrTimestamp     `ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" json:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"` //

	CN                   AttrString                `ldap:"cn" msgpack:"cn,omitempty" json:"cn,omitempty" redis:"cn" redisearch:"text"`                                                                         // RDN in group's context
	DC                   AttrString                `ldap:"dc" msgpack:"dc,omitempty" json:"dc,omitempty" redis:"dc" redisearch:"text,sortable"`                                                                //
	Description          AttrString                `ldap:"description" msgpack:"description,omitempty" json:"description,omitempty" redis:"description" redisearch:"text"`                                     //
	DestinationIndicator AttrDestinationIndicators `ldap:"destinationIndicator" msgpack:"destinationIndicator,omitempty" json:"destinationIndicator,omitempty" redis:"destinationIndicator" redisearch:"text"` //
	DisplayName          AttrString                `ldap:"displayName" msgpack:"displayName,omitempty" json:"displayName,omitempty" redis:"displayName" redisearch:"text,sortable"`                            //
	GIDNumber            AttrIDNumber              `ldap:"gidNumber" msgpack:"gidNumber,omitempty" json:"gidNumber,omitempty" redis:"gidNumber" redisearch:"numeric"`                                          // Primary GIDNumber in user's context (ignore it) and GIDNumber in group's context.
	HomeDirectory        AttrString                `ldap:"homeDirectory" msgpack:"homeDirectory,omitempty" json:"homeDirectory,omitempty" redis:"homeDirectory" redisearch:"text"`                             //
	IPHostNumber         AttrIPHostNumbers         `ldap:"ipHostNumber" msgpack:"ipHostNumber,omitempty" json:"ipHostNumber,omitempty" redis:"ipHostNumber" redisearch:"text,sortable"`                        //
	Mail                 AttrMails                 `ldap:"mail" msgpack:"mail,omitempty" json:"mail,omitempty" redis:"mail" redisearch:"text"`                                                                 //
	Member               AttrDNs                   `ldap:"member" msgpack:"member,omitempty" json:"member,omitempty" redis:"member" redisearch:"tag,sortable"`                                                 //
	MemberOf             AttrDNs                   `ldap:"memberOf" msgpack:"memberOf,omitempty" json:"memberOf,omitempty" redis:"memberOf" redisearch:"tag"`                                                  // ignore it, don't cache, calculate on the fly or avoid
	O                    AttrString                `ldap:"o" msgpack:"o,omitempty" json:"o,omitempty" redis:"o" redisearch:"text"`                                                                             //
	OU                   AttrString                `ldap:"ou" msgpack:"ou,omitempty" json:"ou,omitempty" redis:"ou" redisearch:"text"`                                                                         //
	Owner                AttrDNs                   `ldap:"owner" msgpack:"owner,omitempty" json:"owner,omitempty" redis:"owner" redisearch:"tag"`                                                              //
	SN                   AttrString                `ldap:"sn" msgpack:"sn,omitempty" json:"sn,omitempty" redis:"sn" redisearch:"text"`                                                                         //
	SSHPublicKey         AttrSSHPublicKeys         `ldap:"sshPublicKey" msgpack:"sshPublicKey,omitempty" json:"sshPublicKey,omitempty" redis:"sshPublicKey" redisearch:"tag"`                                  //
	TelephoneNumber      AttrStrings               `ldap:"telephoneNumber" msgpack:"telephoneNumber,omitempty" json:"telephoneNumber,omitempty" redis:"telephoneNumber" redisearch:"text"`                     //
	TelexNumber          AttrStrings               `ldap:"telexNumber" msgpack:"telexNumber,omitempty" json:"telexNumber,omitempty" redis:"telexNumber" redisearch:"text"`                                     //
	UID                  AttrID                    `ldap:"uid" msgpack:"uid,omitempty" json:"uid,omitempty" redis:"uid" redisearch:"text,sortable"`                                                            // RDN in user's context
	UIDNumber            AttrIDNumber              `ldap:"uidNumber" msgpack:"uidNumber,omitempty" json:"uidNumber,omitempty" redis:"uidNumber" redisearch:"numeric,sortable"`                                 //
	UserPKCS12           AttrUserPKCS12s           `ldap:"userPKCS12" msgpack:"userPKCS12,omitempty" json:"userPKCS12,omitempty" redis:"userPKCS12" redisearch:"tag"`                                          //
	UserPassword         AttrUserPassword          `ldap:"userPassword" msgpack:"userPassword,omitempty" json:"userPassword,omitempty" redis:"userPassword" redisearch:"text"`                                 //

	LabelledURI AttrLabeledURIs `ldap:"labelledURI" msgpack:"labelledURI,omitempty" json:"labelledURI,omitempty" redis:"labelledURI" redisearch:"tag"` //
}

// type AttrDN *ldap.DN //

type attrCN string                                     //
type attrCreateTimestamp time.Time                     //
type attrCreatorsName AttrDN                           //
type AttrDN string                                     //
type attrDescription string                            //
type AttrDestinationIndicators []string                // interim host list
type attrDisplayName string                            //
type attrEntryUUID uuid.UUID                           //
type attrGIDNumber uint64                              //
type attrHomeDirectory string                          //
type AttrIPHostNumbers map[netip.Prefix]struct{}       //
type AttrLabeledURIs []LabeledURILegacy                // custom schema alternative TO DO implement custom schemas
type AttrMails []string                                //
type attrMembers []AttrDN                              //
type attrMembersOf []AttrDN                            //
type attrModifiersName AttrDN                          //
type attrModifyTimestamp time.Time                     //
type attrO string                                      //
type attrOU string                                     //
type attrObjectClasses []string                        //
type attrOwners []AttrDN                               //
type attrSN string                                     //
type AttrSSHPublicKeys map[string]mod_ssh.PublicKey    //
type attrTelephoneNumbers []string                     //
type attrTelexNumbers []string                         //
type attrUID string                                    //
type attrUIDNumber uint64                              //
type AttrUserPKCS12s map[AttrDN]mod_crypto.Certificate // any type of cert-key pairs list TODO implement seamless migration from any to P12
type AttrUserPassword string                           //

type AttrDNs []AttrDN           //
type AttrObjectClasses []string //
type AttrID string              //
type AttrIDNumber uint64        //
type AttrString string          //
type AttrStrings []AttrString   //
type AttrTimestamp time.Time    //
type AttrUUID uuid.UUID         //

type LabeledURI struct {
	// XMLName     xml.Name             `xml:"luri"`
	Type        string     `xml:"type,attr,omitempty"` // `(provider|interim|openvpn|ciscovpn)`
	ASN         uint32     `xml:"asn,attr,omitempty"`
	UpstreamASN uint32     `xml:"upstream_asn,attr,omitempty"`
	HostASN     uint32     `xml:"host_asn,attr,omitempty"`
	URL         url.URL    `xml:"url,attr,omitempty"`
	Listen      netip.Addr `xml:"listen,attr,omitempty"`
	ACL         string     `xml:"acl,attr,omitempty"`
	AAA         string     `xml:"aaa,attr,omitempty"`

	OpenVPN     []mod_net.OpenVPN     `xml:"OpenVPN,omitempty"`
	CiscoVPN    []mod_net.CiscoVPN    `xml:"CiscoVPN,omitempty"`
	InterimHost []mod_net.InterimHost `xml:"InterimHost,omitempty"`
	Legacy      []LabeledURILegacy    `xml:"LabelledURI,omitempty"`
}
type LabeledURILegacy struct {
	Key   string `xml:"key,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}
