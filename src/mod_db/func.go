package mod_db

import (
	"context"

	"github.com/RediSearch/redisearch-go/redisearch"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	var (
		docs   []*redisearch.Document
		schema *redisearch.Schema
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

	switch err = outbound.rsClient.CreateIndexWithIndexDefinition(
		schema,
		redisearch.NewIndexDefinition().AddPrefix("ldap:entry:"),
	); {
	case mod_errors.Contains(err, mod_errors.EIndexExist):
	case err != nil:
		return
	}

	for _, doc := range docs {
		switch swErr := outbound.rsClient.Index([]redisearch.Document{*doc}...); {
		case mod_errors.Contains(swErr, mod_errors.EDocExist):
			l.Z{l.E: mod_errors.EDocExist, l.M: doc.Properties[_dn.String()]}.Warning()
		case swErr != nil:
			return swErr
		}
	}

	return
}
