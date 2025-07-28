package mod_db

import (
	"net/netip"
	"time"

	"github.com/google/uuid"

	"rmm23/src/mod_crypto"
	"rmm23/src/mod_ssh"
)

type attrEntryType int   //
type attrEntryStatus int //

type attrCN string                                     //
type attrCreateTimestamp time.Time                     //
type attrCreatorsName attrDN                           //
type attrDN string                                     //
type attrDNs []attrDN                                  //
type attrDescription string                            //
type attrDestinationIndicators []string                // interim host list
type attrDisplayName string                            //
type attrEntryUUID uuid.UUID                           //
type attrGIDNumber uint64                              //
type attrHomeDirectory string                          //
type attrID string                                     //
type attrIDNumber uint64                               //
type attrIPHostNumbers []netip.Prefix                  //
type attrLabeledURIs []labeledURILegacy                //
type attrMails []string                                //
type attrMembers []attrDN                              //
type attrMembersOf []attrDN                            //
type attrModifiersName attrDN                          //
type attrModifyTimestamp time.Time                     //
type attrO string                                      //
type attrOU string                                     //
type attrObjectClasses []string                        //
type attrOwners []attrDN                               //
type attrSN string                                     //
type attrSSHPublicKeys map[string]mod_ssh.PublicKey    //
type attrString string                                 //
type attrStrings []string                              //
type attrTelephoneNumbers []string                     //
type attrTelexNumbers []string                         //
type attrTimestamp time.Time                           //
type attrUID string                                    //
type attrUIDNumber uint64                              //
type attrUUID uuid.UUID                                //
type attrUserPKCS12s map[attrDN]mod_crypto.Certificate // any type of cert-key pairs list (transcoding may apply)
type attrUserPassword string                           //

type labeledURILegacy struct {
	Key   string `json:"key,omitempty"   msgpack:"key,omitempty"`   //
	Value string `json:"value,omitempty" msgpack:"value,omitempty"` //
}
