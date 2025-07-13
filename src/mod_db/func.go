package mod_db

import (
	"context"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/google/uuid"

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

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf) (err error) {
	switch err = inbound.Fetch(); {
	case err != nil:
		return
	}

	var (
		rsClient = redisearch.NewClient("10.133.0.223:6379", "entryIdx")
		entry    = Entry{}
		_        = rsClient.Drop()
	)

	switch err = rsClient.CreateIndex(entry.RedisearchSchema()); {
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
			doc = redisearch.NewDocument(uuid.UUID(interim.UUID).String(), 1.0)
		)
		doc.Set("Type", interim.Type)
		doc.Set("UUID", interim.UUID)
		doc.Set("DN", interim.DN)
		doc.Set("ObjectClass", interim.ObjectClass)
		doc.Set("CreatorsName", interim.CreatorsName)
		doc.Set("CreateTimestamp", interim.CreateTimestamp)
		doc.Set("ModifiersName", interim.ModifiersName)
		doc.Set("ModifyTimestamp", interim.ModifyTimestamp)
		doc.Set("DC", interim.DC)
		doc.Set("O", interim.O)
		doc.Set("Legacy", interim.Legacy)

		switch err = rsClient.Index([]redisearch.Document{doc}...); {
		case err != nil:
			return
		}

		// var (
		// 	jsonData []byte
		// )
		//
		// switch jsonData, err = json.Marshal(interim); {
		// case err != nil:
		// 	return
		// }
		//
		// switch err = rdb.Set(ctx, uuid.UUID(interim.UUID).String(), jsonData, 0).Err(); {
		// case err != nil:
		// 	return
		// }

	}
	//

	return
}
