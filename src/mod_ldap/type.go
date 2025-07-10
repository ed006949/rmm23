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

type Entry struct{ *ldap.Entry }

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

	domain *ElementDomain
	users  ElementUsers
	groups ElementGroups
	hosts  ElementHosts

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

type ElementHosts map[AttrDN]*ElementHost
type ElementUsers map[AttrDN]*ElementUser
type ElementGroups map[AttrDN]*ElementGroup

type ElementDomain struct {
	UUID AttrUUID `ldap:"entryUUID"`

	DN          AttrDN            `ldap:"dn"`
	ObjectClass AttrObjectClasses `ldap:"objectClass"`

	DC AttrString `ldap:"dc"`
	O  AttrString `ldap:"o"`

	CreatorsName    AttrDN        `ldap:"creatorsName"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp"`
	ModifiersName   AttrDN        `ldap:"modifiersName"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp"`
}
type ElementUser struct {
	UUID                 AttrUUID                  `ldap:"entryUUID" redis:"uuid" redisearch:"text,sortable"`
	DN                   AttrDN                    `ldap:"dn" redis:"dn" redisearch:"text,sortable"`
	ObjectClass          AttrObjectClasses         `ldap:"objectClass" redis:"objectClass" redisearch:"text"`
	CN                   AttrString                `ldap:"cn" redis:"cn" redisearch:"text"`
	Description          AttrString                `ldap:"description" redis:"description" redisearch:"text"`
	DestinationIndicator AttrDestinationIndicators `ldap:"destinationIndicator" redis:"destinationIndicator" redisearch:"text"`
	DisplayName          AttrString                `ldap:"displayName" redis:"displayName" redisearch:"text,sortable"`
	GIDNumber            AttrIDNumber              `ldap:"gidNumber" redis:"gidNumber" redisearch:"numeric"` // primary group id number
	HomeDirectory        AttrString                `ldap:"homeDirectory" redis:"homeDirectory" redisearch:"text"`
	IPHostNumber         AttrIPHostNumbers         `ldap:"ipHostNumber" redis:"ipHostNumber" redisearch:"text,sortable"`
	LabeledURI           AttrLabeledURIs           `ldap:"labeledURI"`
	Mail                 AttrMails                 `ldap:"mail" redis:"mail" redisearch:"text"`
	MemberOf             AttrDNs                   `ldap:"memberOf"`
	O                    AttrString                `ldap:"o" redis:"o" redisearch:"text"`
	OU                   AttrString                `ldap:"ou" redis:"ou" redisearch:"text"`
	SN                   AttrString                `ldap:"sn" redis:"sn" redisearch:"text"`
	SSHPublicKey         AttrSSHPublicKeys         `ldap:"sshPublicKey" redis:"sshPublicKey" redisearch:"text"`
	TelephoneNumber      AttrStrings               `ldap:"telephoneNumber" redis:"telephoneNumber" redisearch:"text"`
	TelexNumber          AttrStrings               `ldap:"telexNumber" redis:"telexNumber" redisearch:"text"`
	UID                  AttrID                    `ldap:"uid" redis:"uid" redisearch:"text,sortable"`
	UIDNumber            AttrIDNumber              `ldap:"uidNumber" redis:"uidNumber" redisearch:"numeric,sortable"`
	UserPKCS12           AttrUserPKCS12s           `ldap:"userPKCS12"`
	UserPassword         AttrUserPassword          `ldap:"userPassword" redis:"userPassword" redisearch:"text"`
	CreatorsName         AttrDN                    `ldap:"creatorsName" redis:"creatorsName" redisearch:"text"`
	CreateTimestamp      AttrTimestamp             `ldap:"createTimestamp" redis:"createTimestamp" redisearch:"text"`
	ModifiersName        AttrDN                    `ldap:"modifiersName" redis:"modifiersName" redisearch:"text"`
	ModifyTimestamp      AttrTimestamp             `ldap:"modifyTimestamp" redis:"modifyTimestamp" redisearch:"text"`
}
type ElementGroup struct {
	UUID AttrUUID `ldap:"entryUUID"`

	DN          AttrDN            `ldap:"dn"`
	ObjectClass AttrObjectClasses `ldap:"objectClass"`

	CN         AttrString      `ldap:"cn"`
	GIDNumber  AttrIDNumber    `ldap:"gidNumber"`
	LabeledURI AttrLabeledURIs `ldap:"labeledURI"`
	Member     AttrDNs         `ldap:"member"`
	Owner      AttrDNs         `ldap:"owner"`

	CreatorsName    AttrDN        `ldap:"creatorsName"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp"`
	ModifiersName   AttrDN        `ldap:"modifiersName"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp"`
}
type ElementHost struct {
	UUID AttrUUID `ldap:"entryUUID"`

	DN          AttrDN            `ldap:"dn"`
	ObjectClass AttrObjectClasses `ldap:"objectClass"`

	CN            AttrString      `ldap:"cn"`
	GIDNumber     AttrIDNumber    `ldap:"gidNumber"`
	HomeDirectory AttrString      `ldap:"homeDirectory"`
	LabeledURI    AttrLabeledURIs `ldap:"labeledURI"`
	MemberOf      AttrDNs         `ldap:"memberOf"`
	SN            AttrString      `ldap:"sn"`
	UID           AttrID          `ldap:"uid"`
	UIDNumber     AttrIDNumber    `ldap:"uidNumber"`
	UserPKCS12    AttrUserPKCS12s `ldap:"userPKCS12"`

	CreatorsName    AttrDN        `ldap:"creatorsName"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp"`
	ModifiersName   AttrDN        `ldap:"modifiersName"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp"`
}

// type AttrDN *ldap.DN //

type attrCN string                                 //
type attrCreateTimestamp time.Time                 //
type attrCreatorsName AttrDN                       //
type AttrDN string                                 //
type attrDescription string                        //
type AttrDestinationIndicators map[string]struct{} // interim host list
type attrDisplayName string                        //
type attrEntryUUID uuid.UUID                       //
type attrGIDNumber uint64                          //
type attrHomeDirectory string                      //
type AttrIPHostNumbers struct {
	modified bool
	// invalid  error
	data netip.Prefix
} //
type AttrLabeledURIs struct {
	modified bool
	// invalid  error
	data *LabeledURI
}                                                      // custom schema alternative TO DO implement custom schemas
type AttrMails map[string]struct{}                     //
type attrMembers map[AttrDN]struct{}                   //
type attrMembersOf map[AttrDN]struct{}                 //
type attrModifiersName AttrDN                          //
type attrModifyTimestamp time.Time                     //
type attrO string                                      //
type attrOU string                                     //
type attrObjectClasses map[string]struct{}             //
type attrOwners map[AttrDN]struct{}                    //
type attrSN string                                     //
type AttrSSHPublicKeys map[string]mod_ssh.PublicKey    //
type attrTelephoneNumbers map[string]struct{}          //
type attrTelexNumbers map[string]struct{}              //
type attrUID string                                    //
type attrUIDNumber uint64                              //
type AttrUserPKCS12s map[AttrDN]mod_crypto.Certificate // any type of cert-key pairs list TODO implement seamless migration from any to P12
type AttrUserPassword string                           //

type AttrDNs map[AttrDN]struct{}           //
type AttrObjectClasses map[string]struct{} //
type AttrID string                         //
type AttrIDNumber uint64                   //
type AttrString string                     //
type AttrStrings map[AttrString]struct{}   //
type AttrTimestamp time.Time               //
type AttrUUID uuid.UUID                    //

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
