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

	Domain *ElementDomain
	Users  ElementUsers
	Groups ElementGroups
	Hosts  ElementHosts

	searchResults map[string]*ldap.SearchResult
}
type ConfTable struct {
	Domain       map[AttrDN]*ElementDomain
	Users        map[AttrDN]*ElementUsers
	Groups       map[AttrDN]*ElementGroups
	Hosts        map[AttrDN]*ElementHosts
	IPHostNumber map[netip.Prefix]struct{}
	ID           map[AttrIDNumber]struct{}
}

type Elements map[AttrDN]*Element
type ElementHosts map[AttrDN]*ElementHost
type ElementUsers map[AttrDN]*ElementUser
type ElementGroups map[AttrDN]*ElementGroup

type Element struct {
	UUID            AttrUUID          `ldap:"entryUUID" msgpack:"entryUUID,omitempty" json:"entryUUID,omitempty" redis:"uuid" redisearch:"text,sortable"`                     // must be unique
	DN              AttrDN            `ldap:"dn" msgpack:"dn,omitempty" json:"dn,omitempty" redis:"dn" redisearch:"text,sortable"`                                            // must be unique
	ObjectClass     AttrObjectClasses `ldap:"objectClass" msgpack:"objectClass,omitempty" json:"objectClass,omitempty" redis:"objectClass" redisearch:"tag"`                  // entry type
	CreatorsName    AttrDN            `ldap:"creatorsName" msgpack:"creatorsName,omitempty" json:"creatorsName,omitempty" redis:"creatorsName" redisearch:"text"`             //
	CreateTimestamp AttrTimestamp     `ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" json:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"` //
	ModifiersName   AttrDN            `ldap:"modifiersName" msgpack:"modifiersName,omitempty" json:"modifiersName,omitempty" redis:"modifiersName" redisearch:"text"`         //
	ModifyTimestamp AttrTimestamp     `ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" json:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"` //

	CN                   AttrString                `ldap:"cn,omitempty"`                   //
	DC                   AttrString                `ldap:"dc,omitempty"`                   //
	Description          AttrString                `ldap:"description,omitempty"`          //
	DestinationIndicator AttrDestinationIndicators `ldap:"destinationIndicator,omitempty"` //
	DisplayName          AttrString                `ldap:"displayName,omitempty"`          //
	GIDNumber            AttrIDNumber              `ldap:"gidNumber,omitempty"`            //
	HomeDirectory        AttrString                `ldap:"homeDirectory,omitempty"`        //
	IPHostNumber         AttrIPHostNumbers         `ldap:"ipHostNumber,omitempty"`         //
	Mail                 AttrMails                 `ldap:"mail,omitempty"`                 //
	Member               AttrDNs                   `ldap:"member,omitempty"`               //
	MemberOf             AttrDNs                   `ldap:"memberOf,omitempty"`             //
	O                    AttrString                `ldap:"o,omitempty"`                    //
	OU                   AttrString                `ldap:"ou,omitempty"`                   //
	Owner                AttrDNs                   `ldap:"owner,omitempty"`                //
	SN                   AttrString                `ldap:"sn,omitempty"`                   //
	SSHPublicKey         AttrSSHPublicKeys         `ldap:"sshPublicKey,omitempty"`         //
	TelephoneNumber      AttrStrings               `ldap:"telephoneNumber,omitempty"`      //
	TelexNumber          AttrStrings               `ldap:"telexNumber,omitempty"`          //
	UID                  AttrID                    `ldap:"uid,omitempty"`                  //
	UIDNumber            AttrIDNumber              `ldap:"uidNumber,omitempty"`            //
	UserPKCS12           AttrUserPKCS12s           `ldap:"userPKCS12,omitempty"`           //
	UserPassword         AttrUserPassword          `ldap:"userPassword,omitempty"`         //

	LabeledURI AttrLabeledURIs `ldap:"labeledURI" msgpack:"labeledURI,omitempty" json:"labeledURI,omitempty" redis:"labeledURI" redisearch:"text"` //
}

type ElementDomain struct {
	UUID            AttrUUID          `ldap:"entryUUID" msgpack:"entryUUID,omitempty" json:"entryUUID,omitempty" redis:"uuid" redisearch:"text,sortable"`                     // must be unique
	DN              AttrDN            `ldap:"dn" msgpack:"dn,omitempty" json:"dn,omitempty" redis:"dn" redisearch:"text,sortable"`                                            // must be unique
	ObjectClass     AttrObjectClasses `ldap:"objectClass" msgpack:"objectClass,omitempty" json:"objectClass,omitempty" redis:"objectClass" redisearch:"tag"`                  // entry type
	CreatorsName    AttrDN            `ldap:"creatorsName" msgpack:"creatorsName,omitempty" json:"creatorsName,omitempty" redis:"creatorsName" redisearch:"text"`             //
	CreateTimestamp AttrTimestamp     `ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" json:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"` //
	ModifiersName   AttrDN            `ldap:"modifiersName" msgpack:"modifiersName,omitempty" json:"modifiersName,omitempty" redis:"modifiersName" redisearch:"text"`         //
	ModifyTimestamp AttrTimestamp     `ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" json:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"` //

	DC AttrString `ldap:"dc" msgpack:"dc,omitempty" json:"dc,omitempty" redis:"dc" redisearch:"text,sortable"` //
	O  AttrString `ldap:"o" msgpack:"o,omitempty" json:"o,omitempty" redis:"o" redisearch:"text,sortable"`     //

	LabeledURI AttrLabeledURIs `ldap:"labeledURI" msgpack:"labeledURI,omitempty" json:"labeledURI,omitempty" redis:"labeledURI" redisearch:"text"` //
}
type ElementUser struct {
	UUID            AttrUUID          `ldap:"entryUUID"`
	DN              AttrDN            `ldap:"dn"`
	ObjectClass     AttrObjectClasses `ldap:"objectClass"`
	CreatorsName    AttrDN            `ldap:"creatorsName"`
	CreateTimestamp AttrTimestamp     `ldap:"createTimestamp"`
	ModifiersName   AttrDN            `ldap:"modifiersName"`
	ModifyTimestamp AttrTimestamp     `ldap:"modifyTimestamp"`

	CN                   AttrString                `ldap:"cn"`
	Description          AttrString                `ldap:"description"`
	DestinationIndicator AttrDestinationIndicators `ldap:"destinationIndicator"`
	DisplayName          AttrString                `ldap:"displayName"`
	GIDNumber            AttrIDNumber              `ldap:"gidNumber"`
	HomeDirectory        AttrString                `ldap:"homeDirectory"`
	IPHostNumber         AttrIPHostNumbers         `ldap:"ipHostNumber"`
	Mail                 AttrMails                 `ldap:"mail"`
	MemberOf             AttrDNs                   `ldap:"memberOf"`
	O                    AttrString                `ldap:"o"`
	OU                   AttrString                `ldap:"ou"`
	SN                   AttrString                `ldap:"sn"`
	SSHPublicKey         AttrSSHPublicKeys         `ldap:"sshPublicKey"`
	TelephoneNumber      AttrStrings               `ldap:"telephoneNumber"`
	TelexNumber          AttrStrings               `ldap:"telexNumber"`
	UID                  AttrID                    `ldap:"uid"`
	UIDNumber            AttrIDNumber              `ldap:"uidNumber"`
	UserPKCS12           AttrUserPKCS12s           `ldap:"userPKCS12"`
	UserPassword         AttrUserPassword          `ldap:"userPassword"`

	LabeledURI AttrLabeledURIs `ldap:"labeledURI" msgpack:"labeledURI,omitempty" json:"labeledURI,omitempty" redis:"labeledURI" redisearch:"text"` //
}
type ElementGroup struct {
	UUID            AttrUUID          `ldap:"entryUUID"`
	DN              AttrDN            `ldap:"dn"`
	ObjectClass     AttrObjectClasses `ldap:"objectClass"`
	CreatorsName    AttrDN            `ldap:"creatorsName"`
	CreateTimestamp AttrTimestamp     `ldap:"createTimestamp"`
	ModifiersName   AttrDN            `ldap:"modifiersName"`
	ModifyTimestamp AttrTimestamp     `ldap:"modifyTimestamp"`

	CN        AttrString   `ldap:"cn"`
	GIDNumber AttrIDNumber `ldap:"gidNumber"`
	Member    AttrDNs      `ldap:"member"`
	Owner     AttrDNs      `ldap:"owner"`

	LabeledURI AttrLabeledURIs `ldap:"labeledURI" msgpack:"labeledURI,omitempty" json:"labeledURI,omitempty" redis:"labeledURI" redisearch:"text"` //
}
type ElementHost struct {
	UUID            AttrUUID          `ldap:"entryUUID"`
	DN              AttrDN            `ldap:"dn"`
	ObjectClass     AttrObjectClasses `ldap:"objectClass"`
	CreatorsName    AttrDN            `ldap:"creatorsName"`
	CreateTimestamp AttrTimestamp     `ldap:"createTimestamp"`
	ModifiersName   AttrDN            `ldap:"modifiersName"`
	ModifyTimestamp AttrTimestamp     `ldap:"modifyTimestamp"`

	CN            AttrString      `ldap:"cn"`
	GIDNumber     AttrIDNumber    `ldap:"gidNumber"`
	HomeDirectory AttrString      `ldap:"homeDirectory"`
	MemberOf      AttrDNs         `ldap:"memberOf"`
	SN            AttrString      `ldap:"sn"`
	UID           AttrID          `ldap:"uid"`
	UIDNumber     AttrIDNumber    `ldap:"uidNumber"`
	UserPKCS12    AttrUserPKCS12s `ldap:"userPKCS12"`

	LabeledURI AttrLabeledURIs `ldap:"labeledURI" msgpack:"labeledURI,omitempty" json:"labeledURI,omitempty" redis:"labeledURI" redisearch:"text"` //
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
	Legacy      []LabeledURILegacy    `xml:"Legacy,omitempty"`
}
type LabeledURILegacy struct {
	Key   string `xml:"key,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}
