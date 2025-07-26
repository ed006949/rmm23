package mod_db

import (
	"context"
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/google/uuid"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	var (
		docs   []*redisearch.Document
		schema *redisearch.Schema
		// rdocs      []redisearch.Document
		// rdocscount int
	)

	// predefine schema
	switch schema, err = new(Entry).redisearchSchema(); {
	case err != nil:
		return
	}

	switch docs, err = getLDAPDocs(inbound, schema); {
	case err != nil:
		return
	}

	switch err = outbound.dial(); {
	case err != nil:
		return
	}

	switch err = outbound.rsClient.CreateIndex(schema); {
	case mod_errors.Contains(err, mod_errors.EIndexExist):
	case err != nil:
		return
	}

	for _, doc := range docs {
		var (
			idDoc *redisearch.Document
		)

		switch idDoc, err = outbound.getDoc(doc.Id); {
		case err != nil: // error
			return err
		case idDoc != nil: // Document already exist, hash new `UUID` from `DN`
			var (
				count int
			)

			switch _, count, err = outbound.getDocsByKV(_dn, doc.Properties[_dn.String()]); {
			case err != nil: // error
				return
			case count == 0: // same `DN` not exist, hash new `UUID` from `DN`
				var (
					newUUID = mod_ldap.AttrUUID(uuid.NewSHA1(uuid.Nil, []byte(fmt.Sprint(doc.Properties[_dn.String()]))))
				)

				doc.Id = newUUID.Entry()
				doc.Set(_uuid.String(), newUUID)
			case count == 1: // same `DN` exist, skip
				continue
			case count > 1: // multiple 'DN' exist
				return mod_errors.EUnwilling
			}
		}

		switch err = outbound.rsClient.Index([]redisearch.Document{*doc}...); {
		case mod_errors.Contains(err, mod_errors.EDocExist):
			fmt.Print(doc.Id, "\n")
		case err != nil:
			return
		}
	}

	return
}
