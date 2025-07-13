package mod_db

import (
	"net/netip"
	"net/url"

	"github.com/redis/go-redis/v9"

	"rmm23/src/mod_ldap"
)

// type ConfTable struct {
//	Domain       map[AttrDN]*ElementDomain
//	Users        map[AttrDN]*ElementUsers
//	Groups       map[AttrDN]*ElementGroups
//	Hosts        map[AttrDN]*ElementHosts
//	IPHostNumber map[netip.Prefix]struct{}
//	ID           map[AttrIDNumber]struct{}
// }

// type ConfDomain struct {
//    DN            AttrDN `xml:"dn,attr"`
//    domain        *ElementDomain
//    users         ElementUsers
//    groups        ElementGroups
//    hosts         ElementHosts
//    searchResults map[string]*ldap.SearchResult
// }

func CopyLDAP2DB(inbound *mod_ldap.Conf, rdb *redis.Client) (err error) {
	switch err = inbound.Fetch(); {
	case err != nil:
		return
	}

	//
	for _, b := range inbound.Domain {
		var (
			interim = Entry{
				Type:            EntryTypeDomain,
				UUID:            b.Domain.UUID,
				DN:              b.Domain.DN,
				ObjectClass:     b.Domain.ObjectClass,
				CreatorsName:    b.Domain.CreatorsName,
				CreateTimestamp: b.Domain.CreateTimestamp,
				ModifiersName:   b.Domain.ModifiersName,
				ModifyTimestamp: b.Domain.ModifyTimestamp,
				DC:              b.Domain.DC,
				O:               b.Domain.O,
				Legacy:          b.Domain.LabeledURI,
			}
		)

	}
	//

	return
}
