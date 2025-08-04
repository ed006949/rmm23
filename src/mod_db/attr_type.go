package mod_db

import (
	"crypto/x509/pkix"
	"net/netip"
	"time"

	"github.com/google/uuid"

	"rmm23/src/mod_ssh"
)

type attrEntryType int   //
type attrEntryStatus int //

type attrCN string                                  //
type attrCreateTimestamp struct{ time.Time }        //
type attrCreatorsName *attrDN                       //
type attrDN struct{ pkix.Name }                     //
type attrDNs []*attrDN                              //
type attrDescription string                         //
type attrDestinationIndicators []string             // interim host list
type attrDisplayName string                         //
type attrEntryUUID struct{ uuid.UUID }              //
type attrGIDNumber uint64                           //
type attrHomeDirectory string                       //
type attrID string                                  //
type attrIDNumber uint64                            //
type attrIPHostNumbers []netip.Prefix               //
type attrLabeledURIs map[string]string              //
type attrMails []string                             //
type attrMembers []*attrDN                          //
type attrMembersOf []*attrDN                        //
type attrModifiersName *attrDN                      //
type attrModifyTimestamp struct{ time.Time }        //
type attrO string                                   //
type attrOU string                                  //
type attrObjectClasses []string                     //
type attrOwners []*attrDN                           //
type attrSN string                                  //
type attrSSHPublicKeys map[string]mod_ssh.PublicKey //
type attrString string                              //
type attrStrings []string                           //
type attrTelephoneNumbers []string                  //
type attrTelexNumbers []string                      //
type attrTime struct{ time.Time }                   //
type attrUID string                                 //
type attrUIDNumber uint64                           //
type attrUUID struct{ uuid.UUID }                   //
type attrUserPassword string                        //
