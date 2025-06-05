package mod_ldap

import (
	"net/netip"
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

	schema map[string]*schema
	conn   *ldap.Conn

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
	DN          AttrDN      `ldap:"dn"`
	ObjectClass AttrStrings `ldap:"objectClass"`

	DC AttrString `ldap:"dc"`
	O  AttrString `ldap:"o"`

	EntryUUID       AttrUUID      `ldap:"entryUUID"`
	CreatorsName    AttrDN        `ldap:"creatorsName"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp"`
	ModifiersName   AttrDN        `ldap:"modifiersName"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp"`
}

type ElementUser struct {
	DN          AttrDN      `ldap:"dn"`
	ObjectClass AttrStrings `ldap:"objectClass"`

	CN                   AttrString                `ldap:"cn"`
	Description          AttrString                `ldap:"description"`
	DestinationIndicator AttrDestinationIndicators `ldap:"destinationIndicator"`
	DisplayName          AttrString                `ldap:"displayName"`
	GIDNumber            AttrIDNumber              `ldap:"gidNumber"`
	HomeDirectory        AttrString                `ldap:"homeDirectory"`
	IPHostNumber         AttrIPHostNumbers         `ldap:"ipHostNumber"`
	LabeledURI           AttrLabeledURIs           `ldap:"labeledURI"`
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

	EntryUUID       AttrUUID      `ldap:"entryUUID"`
	CreatorsName    AttrDN        `ldap:"creatorsName"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp"`
	ModifiersName   AttrDN        `ldap:"modifiersName"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp"`
}
type ElementGroup struct {
	DN          AttrDN      `ldap:"dn"`
	ObjectClass AttrStrings `ldap:"objectClass"`

	CN         AttrString      `ldap:"cn"`
	GIDNumber  AttrIDNumber    `ldap:"gidNumber"`
	LabeledURI AttrLabeledURIs `ldap:"labeledURI"`
	Member     AttrDNs         `ldap:"member"`
	Owner      AttrDNs         `ldap:"owner"`

	EntryUUID       AttrUUID      `ldap:"entryUUID"`
	CreatorsName    AttrDN        `ldap:"creatorsName"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp"`
	ModifiersName   AttrDN        `ldap:"modifiersName"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp"`
}
type ElementHost struct {
	DN          AttrDN      `ldap:"dn"`
	ObjectClass AttrStrings `ldap:"objectClass"`

	CN            AttrString      `ldap:"cn"`
	GIDNumber     AttrIDNumber    `ldap:"gidNumber"`
	HomeDirectory AttrString      `ldap:"homeDirectory"`
	LabeledURI    AttrLabeledURIs `ldap:"labeledURI"`
	MemberOf      AttrDNs         `ldap:"memberOf"`
	SN            AttrString      `ldap:"sn"`
	UID           AttrID          `ldap:"uid"`
	UIDNumber     AttrIDNumber    `ldap:"uidNumber"`
	UserPKCS12    AttrUserPKCS12s `ldap:"userPKCS12"`

	EntryUUID       AttrUUID      `ldap:"entryUUID"`
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
	invalid  error
	data     netip.Prefix
} //
type AttrLabeledURIs struct {
	modified bool
	invalid  error
	data     *LabeledURI
}                                                       // custom schema alternative TO DO implement custom schemas
type AttrMails map[string]struct{}                      //
type attrMembers map[AttrDN]struct{}                    //
type attrMembersOf map[AttrDN]struct{}                  //
type attrModifiersName AttrDN                           //
type attrModifyTimestamp time.Time                      //
type attrO string                                       //
type attrOU string                                      //
type attrObjectClasses map[string]struct{}              //
type attrOwners map[AttrDN]struct{}                     //
type attrSN string                                      //
type AttrSSHPublicKeys map[string]mod_ssh.PublicKey     //
type attrTelephoneNumbers map[string]struct{}           //
type attrTelexNumbers map[string]struct{}               //
type attrUID string                                     //
type attrUIDNumber uint64                               //
type AttrUserPKCS12s map[AttrDN]*mod_crypto.Certificate // any type of cert-key pairs list TODO implement seamless migration from any to P12
type AttrUserPassword string                            //

type AttrDNs map[AttrDN]struct{}         //
type AttrID string                       //
type AttrIDNumber uint64                 //
type AttrString string                   //
type AttrStrings map[AttrString]struct{} //
type AttrTimestamp time.Time             //
type AttrUUID uuid.UUID                  //

// type Attr Labeled URI map[string]struct{} // custom schema alternative TO DO implement custom schemas

type schema struct {
	OID           string
	Name          string
	Description   string
	Type          string
	Syntax        string
	SubtypeOf     []string
	MustContain   []string
	MayContain    []string
	SingleValue   bool
	Collective    bool
	NoUserMod     bool
	Usage         string
	Obsolete      bool
	SUP           []string
	Structural    bool
	Auxiliary     bool
	Abstract      bool
	RawDefinition string
}

type LabeledURI struct {
	// XMLName     xml.Name             `xml:"luri"`
	OpenVPN     []mod_net.OpenVPN     `xml:"OpenVPN,omitempty"`
	CiscoVPN    []mod_net.CiscoVPN    `xml:"CiscoVPN,omitempty"`
	InterimHost []mod_net.InterimHost `xml:"InterimHost,omitempty"`
	Legacy      []LabeledURILegacy    `xml:"Legacy,omitempty"`
}
type LabeledURILegacy struct {
	Key   string `xml:"key,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
}
