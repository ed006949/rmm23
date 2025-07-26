package mod_db

import (
	"context"
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"

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
		switch err = outbound.rsClient.Index([]redisearch.Document{*doc}...); {
		case mod_errors.Contains(err, mod_errors.EDocExist):
			fmt.Print(doc.Id, "\n")
		case err != nil:
			return
		}
	}

	return
}
