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

type attrDN struct{ pkix.Name }                     //
type attrDNs []*attrDN                              //
type attrID string                                  //
type attrIDNumber uint64                            //
type attrIPAddress struct{ netip.Addr }             //
type attrIPAddresses []*attrIPAddress               //
type attrIPPrefix struct{ netip.Prefix }            //
type attrIPPrefixes []*attrIPPrefix                 //
type attrLabeledURIs map[string]string              //
type attrMails []string                             //
type attrSSHPublicKey mod_ssh.PublicKey             //
type attrSSHPublicKeys map[string]mod_ssh.PublicKey //
type attrString string                              //
type attrStrings []attrString                       //
type attrTime struct{ time.Time }                   //
type attrUUID struct{ uuid.UUID }                   //
type attrUserPassword string                        //

// type attrCN string                                  //
// type attrCreateTimestamp struct{ time.Time }        //
// type attrCreatorsName *attrDN                       //
// type attrDescription string                         //
// type attrDestinationIndicators []string             // interim host list
// type attrDisplayName string                         //
// type attrEntryUUID struct{ uuid.UUID }              //
// type attrGIDNumber uint64                           //
// type attrHomeDirectory string                       //
// type attrMembers attrDNs                            //
// type attrMembersOf attrDNs                          //
// type attrModifiersName *attrDN                      //
// type attrModifyTimestamp struct{ time.Time }        //
// type attrO string                                   //
// type attrOU string                                  //
// type attrObjectClasses []string                     //
// type attrOwners attrDNs                             //
// type attrSN string                                  //
// type attrTelephoneNumbers []string                  //
// type attrTelexNumbers []string                      //
// type attrUID string                                 //
// type attrUIDNumber uint64                           //
