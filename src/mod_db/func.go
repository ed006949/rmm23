package mod_db

import (
	"context"
	"strings"

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
		// _        = rsClient.DropIndex(false)
	)

	switch err = rsClient.CreateIndex(entry.RedisearchSchema()); {
	case err != nil:
		return
	}

	for _, b := range inbound.Domain {
		var (
			doc = redisearch.NewDocument("ldap:entry:"+uuid.UUID(b.Domain.UUID).String(), 1.0)
		)
		doc.Set("Type", EntryTypeDomain)
		doc.Set("UUID", b.Domain.UUID)
		doc.Set("DN", b.Domain.DN)
		doc.Set("ObjectClass", b.Domain.ObjectClass)
		doc.Set("CreatorsName", b.Domain.CreatorsName)
		doc.Set("CreateTimestamp", b.Domain.CreateTimestamp)
		doc.Set("ModifiersName", b.Domain.ModifiersName)
		doc.Set("ModifyTimestamp", b.Domain.ModifyTimestamp)
		doc.Set("DC", b.Domain.DC)
		doc.Set("O", b.Domain.O)
		doc.Set("Legacy", b.Domain.LabeledURI)
		doc.SetPayload(nil)

		switch err = rsClient.Index([]redisearch.Document{doc}...); {
		case err != nil && strings.Contains(err.Error(), "Document already exists"):
		case err != nil:
			return
		}

	}

	return nil
}
