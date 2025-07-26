package mod_db

import (
	"context"
	"errors"
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-ldap/ldap/v3"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_slices"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	var (
		docs       []*redisearch.Document
		schema     *redisearch.Schema
		rdocs      []redisearch.Document
		rdocscount int
	)

	// predefine schema
	switch schema, err = new(Entry).redisearchSchema(); {
	case err != nil:
		return
	}

	var (
		ldap2doc = func(fnBaseDN string, fnSearchResultType string, fnSearchResult *ldap.SearchResult) (fnErr error) {
			for _, fnB := range fnSearchResult.Entries {
				var (
					fnDoc   *redisearch.Document
					fnEntry = new(Entry)
				)

				switch fnErr = mod_ldap.UnmarshalEntry(fnB, fnEntry); {
				case fnErr != nil:
					return
				}

				switch fnErr = fnEntry.Type.Parse(fnSearchResultType); {
				case fnErr != nil:
					return
				}

				fnEntry.BaseDN = mod_ldap.AttrDN(fnBaseDN)

				switch fnDoc, fnErr = newRedisearchDocument(
					schema,
					mod_slices.JoinStrings([]string{entryDocIDHeader, fnEntry.UUID.String()}, ":", mod_slices.FlagNone),
					1.0,
					fnEntry,
					false,
				); {
				case fnErr != nil:
					return
				}

				docs = append(docs, fnDoc)
			}

			return
		}
	)

	switch err = inbound.SearchFn(ldap2doc); {
	case err != nil:
		return
	}

	switch err = outbound.Dial(); {
	case err != nil:
		return
	}

	var (
		rsQuery = redisearch.NewQuery("*").SetReturnFields("uuid", "dn").Limit(0, connMaxPaging)
	)

	switch err = outbound.rsClient.CreateIndex(schema); {
	case mod_errors.Contains(err, mod_errors.EIndexExist):
	case err != nil:
		return
	}

	switch rdocs, rdocscount, err = outbound.rsClient.Search(rsQuery); {
	case err != nil:
		return
	case rdocscount >= connMaxPaging:
		return errors.New("max paging limit reached")
	}

	fmt.Printf("%v", len(rdocs))

	for _, doc := range docs {
		switch err = outbound.rsClient.Index([]redisearch.Document{*doc}...); {
		case mod_errors.Contains(err, mod_errors.EDocExist):
		case err != nil:
			return
		}
	}

	return
}
