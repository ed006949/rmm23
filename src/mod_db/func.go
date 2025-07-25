package mod_db

import (
	"context"

	"github.com/RediSearch/redisearch-go/redisearch"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_slices"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	// switch err = inbound.Search(); {
	// case err != nil:
	// 	return
	// }
	var (
		docs   []*redisearch.Document
		schema = new(Entry).redisearchSchema() // predefine schema
	)

	switch err = ldap2doc(inbound, schema, docs); {
	case err != nil:
		return
	}

	switch err = outbound.New(); {
	case err != nil:
		return
	}

	var (
		rsQuery = redisearch.NewQuery("*").SetReturnFields("uuid", "dn").Limit(0, connMaxPaging)
	)

	switch err = outbound.rsClient.CreateIndex(schema); {
	// case mod_errors.Contains(err, mod_errors.EIndexExist):
	case err != nil:
		return
	}

	switch a, b, c := outbound.rsClient.Search(rsQuery); {
	case c != nil:
		return c
	default:
		a = a
		b = b

		panic(nil)
	}

	for _, doc := range docs {
		switch err = outbound.rsClient.Index([]redisearch.Document{*doc}...); {
		case mod_errors.Contains(err, mod_errors.EDocExist):
		case err != nil:
			return
		}
	}

	return
}
func ldap2doc(inbound *mod_ldap.Conf, schema *redisearch.Schema, docs []*redisearch.Document) (err error) {
	for _, b := range inbound.Domains {
		for c, d := range b.SearchResults {
			var (
				entryType AttrType
			)
			switch err = entryType.Parse(c); {
			case err != nil:
				return
			}

			for _, f := range d.Entries {
				var (
					doc   *redisearch.Document
					entry = new(Entry)
				)

				switch err = mod_ldap.UnmarshalEntry(f, entry); {
				case err != nil:
					return
				}

				entry.Type = entryType
				entry.BaseDN = b.DN

				switch doc, err = newRedisearchDocument(
					schema,
					mod_slices.JoinStrings([]string{entryDocIDHeader, entry.UUID.String()}, ":", mod_slices.FlagNone),
					1.0,
					entry,
					false,
				); {
				case err != nil:
					return
				}

				docs = append(docs, doc)
			}
		}
	}

	return
}
