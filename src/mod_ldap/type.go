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

type Conf struct {
	URL      *mod_net.URL `json:"url,omitempty"`
	Settings []*settings  `json:"settings,omitempty"`
	Domains  []*domain    `json:"domain,omitempty"`
	conn     *ldap.Conn
}

type settings struct {
	Type   string           `json:"type,omitempty"`
	DN     AttrDN           `json:"dn,omitempty"`
	CN     string           `json:"cn,omitempty"`
	Scope  AttrSearchScope  `json:"scope,omitempty"`
	Filter AttrSearchFilter `json:"filter,omitempty"`
}

type domain struct {
	DN AttrDN `json:"dn,omitempty"`

	SearchResults map[string]*ldap.SearchResult
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
type AttrIPHostNumbers []netip.Prefix                  //
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

type LabeledURILegacy struct {
	Key   string
	Value string
}

type AttrSearchScope int
type AttrSearchFilter string
