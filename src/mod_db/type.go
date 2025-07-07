package mod_db

import (
	"net/netip"
	"net/url"
	"time"

	"github.com/google/uuid"

	"rmm23/src/mod_crypto"
	"rmm23/src/mod_net"
	"rmm23/src/mod_ssh"
)

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

type AttrDNs map[AttrDN]struct{}         //
type AttrID string                       //
type AttrIDNumber uint64                 //
type AttrString string                   //
type AttrStrings map[AttrString]struct{} //
type AttrTimestamp time.Time             //
type AttrUUID uuid.UUID                  //

type ElementEntry struct {
	// element meta
	UUID            AttrUUID      `ldap:"entryUUID" msgpack:"entryUUID" redis:"uuid" redisearch:"text,sortable"`
	DN              AttrDN        `ldap:"dn" msgpack:"dn" redis:"dn" redisearch:"text,sortable"`
	ObjectClass     AttrStrings   `ldap:"objectClass" msgpack:"objectClass" redis:"objectClass" redisearch:"text"`
	CreatorsName    AttrDN        `ldap:"creatorsName" msgpack:"creatorsName" redis:"creatorsName" redisearch:"text"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp" msgpack:"createTimestamp" redis:"createTimestamp" redisearch:"text"`
	ModifiersName   AttrDN        `ldap:"modifiersName" msgpack:"modifiersName" redis:"modifiersName" redisearch:"text"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp" msgpack:"modifyTimestamp" redis:"modifyTimestamp" redisearch:"text"`

	// element data
	CN                   AttrString                `ldap:"cn" msgpack:"cn" redis:"cn" redisearch:"text"`
	DC                   AttrString                `ldap:"dc" msgpack:"dc" redis:"dc" redisearch:"text,sortable"`
	Description          AttrString                `ldap:"description" msgpack:"description" redis:"description" redisearch:"text"`
	DestinationIndicator AttrDestinationIndicators `ldap:"destinationIndicator" msgpack:"destinationIndicator" redis:"destinationIndicator" redisearch:"text"`
	DisplayName          AttrString                `ldap:"displayName" msgpack:"displayName" redis:"displayName" redisearch:"text,sortable"`
	GIDNumber            AttrIDNumber              `ldap:"gidNumber" msgpack:"gidNumber" redis:"gidNumber" redisearch:"numeric"`
	HomeDirectory        AttrString                `ldap:"homeDirectory" msgpack:"homeDirectory" redis:"homeDirectory" redisearch:"text"`
	IPHostNumber         AttrIPHostNumbers         `ldap:"ipHostNumber" msgpack:"ipHostNumber" redis:"ipHostNumber" redisearch:"text,sortable"`
	Mail                 AttrMails                 `ldap:"mail" msgpack:"mail" redis:"mail" redisearch:"text"`
	Member               AttrDNs                   `ldap:"member" msgpack:"member" redis:"member" redisearch:"text,sortable"`
	MemberOf             AttrDNs                   `ldap:"memberOf" msgpack:"memberOf" redis:"memberOf" redisearch:"text"`
	O                    AttrString                `ldap:"o" msgpack:"o" redis:"o" redisearch:"text"`
	OU                   AttrString                `ldap:"ou" msgpack:"ou" redis:"ou" redisearch:"text"`
	Owner                AttrDNs                   `ldap:"owner" msgpack:"owner" redis:"owner" redisearch:"text"`
	SN                   AttrString                `ldap:"sn" msgpack:"sn" redis:"sn" redisearch:"text"`
	SSHPublicKey         AttrSSHPublicKeys         `ldap:"sshPublicKey" msgpack:"sshPublicKey" redis:"sshPublicKey" redisearch:"text"`
	TelephoneNumber      AttrStrings               `ldap:"telephoneNumber" msgpack:"telephoneNumber" redis:"telephoneNumber" redisearch:"text"`
	TelexNumber          AttrStrings               `ldap:"telexNumber" msgpack:"telexNumber" redis:"telexNumber" redisearch:"text"`
	UID                  AttrID                    `ldap:"uid" msgpack:"uid" redis:"uid" redisearch:"text,sortable"`
	UIDNumber            AttrIDNumber              `ldap:"uidNumber" msgpack:"uidNumber" redis:"uidNumber" redisearch:"numeric,sortable"`
	UserPKCS12           AttrUserPKCS12s           `ldap:"userPKCS12" msgpack:"userPKCS12" redis:"userPKCS12" redisearch:"text"`
	UserPassword         AttrUserPassword          `ldap:"userPassword" msgpack:"userPassword" redis:"userPassword" redisearch:"text"`

	// host specific data
	Type        string     `xml:"type,attr,omitempty" msgpack:"type" redis:"host_type" redisearch:"text"` // host type `(provider|interim|openvpn|ciscovpn)`
	ASN         uint32     `xml:"asn,attr,omitempty" msgpack:"asn" redis:"host_asn" redisearch:"numeric,sortable"`
	UpstreamASN uint32     `xml:"upstream_asn,attr,omitempty" msgpack:"upstream_asn" redis:"host_upstream_asn" redisearch:"numeric"`
	HostUUID    uint32     `xml:"host_uuid,attr,omitempty" msgpack:"host_uuid" redis:"host_host_uuid" redisearch:"numeric"` // (?) replace with member/memberOf
	URL         url.URL    `xml:"url,attr,omitempty"`
	Listen      netip.Addr `xml:"listen,attr,omitempty"`
	ACL         string     `xml:"acl,attr,omitempty"`
	AAA         string     `xml:"aaa,attr,omitempty"`
}

type ElementDomain struct {
	UUID            AttrUUID      `ldap:"entryUUID" msgpack:"entryUUID" redis:"uuid" redisearch:"text,sortable"`
	DN              AttrDN        `ldap:"dn" msgpack:"dn" redis:"dn" redisearch:"text,sortable"`
	ObjectClass     AttrStrings   `ldap:"objectClass" msgpack:"objectClass" redis:"objectClass" redisearch:"text"`
	CreatorsName    AttrDN        `ldap:"creatorsName" msgpack:"creatorsName" redis:"creatorsName" redisearch:"text"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp" msgpack:"createTimestamp" redis:"createTimestamp" redisearch:"text"`
	ModifiersName   AttrDN        `ldap:"modifiersName" msgpack:"modifiersName" redis:"modifiersName" redisearch:"text"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp" msgpack:"modifyTimestamp" redis:"modifyTimestamp" redisearch:"text"`

	DC AttrString `ldap:"dc" msgpack:"dc" redis:"dc" redisearch:"text,sortable"`
	O  AttrString `ldap:"o" msgpack:"o" redis:"o" redisearch:"text,sortable"`
}
type ElementGroup struct {
	UUID            AttrUUID      `ldap:"entryUUID" msgpack:"entryUUID" redis:"uuid" redisearch:"text,sortable"`
	DN              AttrDN        `ldap:"dn" msgpack:"dn" redis:"dn" redisearch:"text,sortable"`
	ObjectClass     AttrStrings   `ldap:"objectClass" msgpack:"objectClass" redis:"objectClass" redisearch:"text"`
	CreatorsName    AttrDN        `ldap:"creatorsName" msgpack:"creatorsName" redis:"creatorsName" redisearch:"text"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp" msgpack:"createTimestamp" redis:"createTimestamp" redisearch:"text"`
	ModifiersName   AttrDN        `ldap:"modifiersName" msgpack:"modifiersName" redis:"modifiersName" redisearch:"text"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp" msgpack:"modifyTimestamp" redis:"modifyTimestamp" redisearch:"text"`

	CN        AttrString   `ldap:"cn" msgpack:"cn" redis:"cn" redisearch:"text"`
	GIDNumber AttrIDNumber `ldap:"gidNumber" msgpack:"gidNumber" redis:"gidNumber" redisearch:"numeric,sortable"`
	Member    AttrDNs      `ldap:"member" msgpack:"member" redis:"member" redisearch:"text,sortable"`
	Owner     AttrDNs      `ldap:"owner" msgpack:"owner" redis:"owner" redisearch:"text"`

	LabeledURI AttrLabeledURIs `ldap:"labeledURI"`
}
type ElementUser struct {
	UUID            AttrUUID      `ldap:"entryUUID" msgpack:"entryUUID" redis:"uuid" redisearch:"text,sortable"`
	DN              AttrDN        `ldap:"dn" msgpack:"dn" redis:"dn" redisearch:"text,sortable"`
	ObjectClass     AttrStrings   `ldap:"objectClass" msgpack:"objectClass" redis:"objectClass" redisearch:"text"`
	CreatorsName    AttrDN        `ldap:"creatorsName" msgpack:"creatorsName" redis:"creatorsName" redisearch:"text"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp" msgpack:"createTimestamp" redis:"createTimestamp" redisearch:"text"`
	ModifiersName   AttrDN        `ldap:"modifiersName" msgpack:"modifiersName" redis:"modifiersName" redisearch:"text"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp" msgpack:"modifyTimestamp" redis:"modifyTimestamp" redisearch:"text"`

	CN                   AttrString                `ldap:"cn" msgpack:"cn" redis:"cn" redisearch:"text"`
	Description          AttrString                `ldap:"description" msgpack:"description" redis:"description" redisearch:"text"`
	DestinationIndicator AttrDestinationIndicators `ldap:"destinationIndicator" msgpack:"destinationIndicator" redis:"destinationIndicator" redisearch:"text"`
	DisplayName          AttrString                `ldap:"displayName" msgpack:"displayName" redis:"displayName" redisearch:"text,sortable"`
	GIDNumber            AttrIDNumber              `ldap:"gidNumber" msgpack:"gidNumber" redis:"gidNumber" redisearch:"numeric"`
	HomeDirectory        AttrString                `ldap:"homeDirectory" msgpack:"homeDirectory" redis:"homeDirectory" redisearch:"text"`
	IPHostNumber         AttrIPHostNumbers         `ldap:"ipHostNumber" msgpack:"ipHostNumber" redis:"ipHostNumber" redisearch:"text,sortable"`
	Mail                 AttrMails                 `ldap:"mail" msgpack:"mail" redis:"mail" redisearch:"text"`
	MemberOf             AttrDNs                   `ldap:"memberOf" msgpack:"memberOf"`
	O                    AttrString                `ldap:"o" msgpack:"o" redis:"o" redisearch:"text"`
	OU                   AttrString                `ldap:"ou" msgpack:"ou" redis:"ou" redisearch:"text"`
	SN                   AttrString                `ldap:"sn" msgpack:"sn" redis:"sn" redisearch:"text"`
	SSHPublicKey         AttrSSHPublicKeys         `ldap:"sshPublicKey" msgpack:"sshPublicKey" redis:"sshPublicKey" redisearch:"text"`
	TelephoneNumber      AttrStrings               `ldap:"telephoneNumber" msgpack:"telephoneNumber" redis:"telephoneNumber" redisearch:"text"`
	TelexNumber          AttrStrings               `ldap:"telexNumber" msgpack:"telexNumber" redis:"telexNumber" redisearch:"text"`
	UID                  AttrID                    `ldap:"uid" msgpack:"uid" redis:"uid" redisearch:"text,sortable"`
	UIDNumber            AttrIDNumber              `ldap:"uidNumber" msgpack:"uidNumber" redis:"uidNumber" redisearch:"numeric,sortable"`
	UserPKCS12           AttrUserPKCS12s           `ldap:"userPKCS12" msgpack:"userPKCS12" redis:"userPKCS12" redisearch:"text"`
	UserPassword         AttrUserPassword          `ldap:"userPassword" msgpack:"userPassword" redis:"userPassword" redisearch:"text"`

	LabeledURI AttrLabeledURIs `ldap:"labeledURI"`
}
type ElementHost struct {
	UUID            AttrUUID      `ldap:"entryUUID" msgpack:"entryUUID" redis:"uuid" redisearch:"text,sortable"`
	DN              AttrDN        `ldap:"dn" msgpack:"dn" redis:"dn" redisearch:"text,sortable"`
	ObjectClass     AttrStrings   `ldap:"objectClass" msgpack:"objectClass" redis:"objectClass" redisearch:"text"`
	CreatorsName    AttrDN        `ldap:"creatorsName" msgpack:"creatorsName" redis:"creatorsName" redisearch:"text"`
	CreateTimestamp AttrTimestamp `ldap:"createTimestamp" msgpack:"createTimestamp" redis:"createTimestamp" redisearch:"text"`
	ModifiersName   AttrDN        `ldap:"modifiersName" msgpack:"modifiersName" redis:"modifiersName" redisearch:"text"`
	ModifyTimestamp AttrTimestamp `ldap:"modifyTimestamp" msgpack:"modifyTimestamp" redis:"modifyTimestamp" redisearch:"text"`

	CN            AttrString   `ldap:"cn" msgpack:"cn" redis:"cn" redisearch:"text"`
	GIDNumber     AttrIDNumber `ldap:"gidNumber" msgpack:"gidNumber" redis:"gidNumber" redisearch:"numeric"`
	HomeDirectory AttrString   `ldap:"homeDirectory" msgpack:"homeDirectory" redis:"homeDirectory" redisearch:"text"`
	// MemberOf      AttrDNs         `ldap:"memberOf"`
	SN         AttrString      `ldap:"sn" msgpack:"sn" redis:"sn" redisearch:"text"`
	UID        AttrID          `ldap:"uid" msgpack:"uid" redis:"uid" redisearch:"text,sortable"`
	UIDNumber  AttrIDNumber    `ldap:"uidNumber" msgpack:"uidNumber" redis:"uidNumber" redisearch:"numeric,sortable"`
	UserPKCS12 AttrUserPKCS12s `ldap:"userPKCS12" msgpack:"userPKCS12" redis:"userPKCS12" redisearch:"text"`

	LabeledURI AttrLabeledURIs `ldap:"labeledURI"`

	// Type `(provider|interim|openvpn|ciscovpn)`
	Type        string     `xml:"type,attr,omitempty" msgpack:"type" redis:"host_type" redisearch:"text"`
	ASN         uint32     `xml:"asn,attr,omitempty" msgpack:"asn" redis:"host_asn" redisearch:"numeric,sortable"`
	UpstreamASN uint32     `xml:"upstream_asn,attr,omitempty" msgpack:"upstream_asn" redis:"host_upstream_asn" redisearch:"numeric"`
	HostUUID    uint32     `xml:"host_uuid,attr,omitempty" msgpack:"host_uuid" redis:"host_host_uuid" redisearch:"numeric"`
	URL         url.URL    `xml:"url,attr,omitempty"`
	Listen      netip.Addr `xml:"listen,attr,omitempty"`
	ACL         string     `xml:"acl,attr,omitempty"`
	AAA         string     `xml:"aaa,attr,omitempty"`

	// (?)
	Member   AttrDNs `ldap:"member" msgpack:"member" redis:"member" redisearch:"text,sortable"`
	Owner    AttrDNs `ldap:"owner" msgpack:"owner" redis:"owner" redisearch:"text"`
	MemberOf AttrDNs `ldap:"memberOf" msgpack:"memberOf" redis:"memberOf" redisearch:"text"`
}

type LabeledURI struct {
	// XMLName     xml.Name             `xml:"luri"`
	Type        string     `xml:"type,attr,omitempty"`
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
