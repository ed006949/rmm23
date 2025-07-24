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

type LDAPConfig struct {
	URL      *mod_net.URL  `json:"url"`
	Settings []LDAPSetting `json:"settings"`
	Domains  []LDAPDomain  `json:"domain"`
	conn     *ldap.Conn

	searchResults map[string]*ldap.SearchResult
}

type LDAPSetting struct {
	Type   string `json:"type"`
	DN     AttrDN `json:"dn"`
	CN     string `json:"cn"`
	Filter string `json:"filter"`

	searchResults map[string]*ldap.SearchResult
}

type LDAPDomain struct {
	DN AttrDN `json:"dn"`

	Domain *Element
	Users  Elements
	Groups Elements
	Hosts  Elements

	searchResults map[string]*ldap.SearchResult
}

type Conf struct {
	URL      *mod_net.URL
	Settings []*ConfSettings
	Domain   []*ConfDomain

	// schema map[string]*schema
	conn *ldap.Conn

	Table *ConfTable
}
type ConfSettings struct {
	Type   string
	DN     AttrDN
	CN     string
	Filter string
}
type ConfDomain struct {
	DN AttrDN

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
	UUID            AttrUUID          `json:"entryUUID,omitempty"       ldap:"entryUUID"       msgpack:"entryUUID,omitempty"       redis:"uuid"            redisearch:"text,sortable"` // must be unique
	DN              AttrDN            `json:"dn,omitempty"              ldap:"dn"              msgpack:"dn,omitempty"              redis:"dn"              redisearch:"text,sortable"` // must be unique
	ObjectClass     AttrObjectClasses `json:"objectClass,omitempty"     ldap:"objectClass"     msgpack:"objectClass,omitempty"     redis:"objectClass"     redisearch:"tag"`           // entry type
	CreatorsName    AttrDN            `json:"creatorsName,omitempty"    ldap:"creatorsName"    msgpack:"creatorsName,omitempty"    redis:"creatorsName"    redisearch:"text"`          //
	CreateTimestamp AttrTimestamp     `json:"createTimestamp,omitempty" ldap:"createTimestamp" msgpack:"createTimestamp,omitempty" redis:"createTimestamp" redisearch:"text"`          //
	ModifiersName   AttrDN            `json:"modifiersName,omitempty"   ldap:"modifiersName"   msgpack:"modifiersName,omitempty"   redis:"modifiersName"   redisearch:"text"`          //
	ModifyTimestamp AttrTimestamp     `json:"modifyTimestamp,omitempty" ldap:"modifyTimestamp" msgpack:"modifyTimestamp,omitempty" redis:"modifyTimestamp" redisearch:"text"`          //

	CN                   AttrString                `json:"cn,omitempty"                   ldap:"cn"                   msgpack:"cn,omitempty"                   redis:"cn"                   redisearch:"text"`             // RDN in group's context
	DC                   AttrString                `json:"dc,omitempty"                   ldap:"dc"                   msgpack:"dc,omitempty"                   redis:"dc"                   redisearch:"text,sortable"`    //
	Description          AttrString                `json:"description,omitempty"          ldap:"description"          msgpack:"description,omitempty"          redis:"description"          redisearch:"text"`             //
	DestinationIndicator AttrDestinationIndicators `json:"destinationIndicator,omitempty" ldap:"destinationIndicator" msgpack:"destinationIndicator,omitempty" redis:"destinationIndicator" redisearch:"text"`             //
	DisplayName          AttrString                `json:"displayName,omitempty"          ldap:"displayName"          msgpack:"displayName,omitempty"          redis:"displayName"          redisearch:"text,sortable"`    //
	GIDNumber            AttrIDNumber              `json:"gidNumber,omitempty"            ldap:"gidNumber"            msgpack:"gidNumber,omitempty"            redis:"gidNumber"            redisearch:"numeric"`          // Primary GIDNumber in user's context (ignore it) and GIDNumber in group's context.
	HomeDirectory        AttrString                `json:"homeDirectory,omitempty"        ldap:"homeDirectory"        msgpack:"homeDirectory,omitempty"        redis:"homeDirectory"        redisearch:"text"`             //
	IPHostNumber         AttrIPHostNumbers         `json:"ipHostNumber,omitempty"         ldap:"ipHostNumber"         msgpack:"ipHostNumber,omitempty"         redis:"ipHostNumber"         redisearch:"text,sortable"`    //
	Mail                 AttrMails                 `json:"mail,omitempty"                 ldap:"mail"                 msgpack:"mail,omitempty"                 redis:"mail"                 redisearch:"text"`             //
	Member               AttrDNs                   `json:"member,omitempty"               ldap:"member"               msgpack:"member,omitempty"               redis:"member"               redisearch:"tag,sortable"`     //
	MemberOf             AttrDNs                   `json:"memberOf,omitempty"             ldap:"memberOf"             msgpack:"memberOf,omitempty"             redis:"memberOf"             redisearch:"tag"`              // ignore it, don't cache, calculate on the fly or avoid
	O                    AttrString                `json:"o,omitempty"                    ldap:"o"                    msgpack:"o,omitempty"                    redis:"o"                    redisearch:"text"`             //
	OU                   AttrString                `json:"ou,omitempty"                   ldap:"ou"                   msgpack:"ou,omitempty"                   redis:"ou"                   redisearch:"text"`             //
	Owner                AttrDNs                   `json:"owner,omitempty"                ldap:"owner"                msgpack:"owner,omitempty"                redis:"owner"                redisearch:"tag"`              //
	SN                   AttrString                `json:"sn,omitempty"                   ldap:"sn"                   msgpack:"sn,omitempty"                   redis:"sn"                   redisearch:"text"`             //
	SSHPublicKey         AttrSSHPublicKeys         `json:"sshPublicKey,omitempty"         ldap:"sshPublicKey"         msgpack:"sshPublicKey,omitempty"         redis:"sshPublicKey"         redisearch:"tag"`              //
	TelephoneNumber      AttrStrings               `json:"telephoneNumber,omitempty"      ldap:"telephoneNumber"      msgpack:"telephoneNumber,omitempty"      redis:"telephoneNumber"      redisearch:"text"`             //
	TelexNumber          AttrStrings               `json:"telexNumber,omitempty"          ldap:"telexNumber"          msgpack:"telexNumber,omitempty"          redis:"telexNumber"          redisearch:"text"`             //
	UID                  AttrID                    `json:"uid,omitempty"                  ldap:"uid"                  msgpack:"uid,omitempty"                  redis:"uid"                  redisearch:"text,sortable"`    // RDN in user's context
	UIDNumber            AttrIDNumber              `json:"uidNumber,omitempty"            ldap:"uidNumber"            msgpack:"uidNumber,omitempty"            redis:"uidNumber"            redisearch:"numeric,sortable"` //
	UserPKCS12           AttrUserPKCS12s           `json:"userPKCS12,omitempty"           ldap:"userPKCS12"           msgpack:"userPKCS12,omitempty"           redis:"userPKCS12"           redisearch:"tag"`              //
	UserPassword         AttrUserPassword          `json:"userPassword,omitempty"         ldap:"userPassword"         msgpack:"userPassword,omitempty"         redis:"userPassword"         redisearch:"text"`             //

	LabeledURI AttrLabeledURIs `json:"labeledURI,omitempty" ldap:"labeledURI" msgpack:"labeledURI,omitempty" redis:"labeledURI" redisearch:"tag"` //
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
type AttrUserPKCS12s map[AttrDN]mod_crypto.Certificate // any type of cert-key pairs list (transcoding may apply)
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
	// XMLName     xml.Name
	Type        string      // `(provider|interim|openvpn|ciscovpn)`
	ASN         uint32
	UpstreamASN uint32
	HostASN     uint32
	URL         url.URL
	Listen      netip.Addr
	ACL         string
	AAA         string

	OpenVPN     []mod_net.OpenVPN
	CiscoVPN    []mod_net.CiscoVPN
	InterimHost []mod_net.InterimHost
	Legacy      []LabeledURILegacy
}
type LabeledURILegacy struct {
	Key   string
	Value string
}
