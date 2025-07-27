package mod_db

import (
	"context"
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	var (
		docs      []*redisearch.Document
		indexInfo *redisearch.IndexInfo
	)

	// define schema
	switch outbound.schema, outbound.schemaMap, err = new(entry).redisearchSchema(); {
	case err != nil:
		return
	}

	switch docs, err = getLDAPDocs(inbound, outbound.schema); {
	case err != nil:
		return
	}

	switch err = outbound.dial(); {
	case err != nil:
		return
	}

	switch err = outbound.createIndex(); {
	case err != nil:
		return
	}

	switch indexInfo, err = outbound.rsClient.Info(); {
	case err != nil:
		return
	}

	l.Z{l.M: indexInfo.Name, "indexing": indexInfo.IsIndexing, "schema": indexInfo.Schema.Fields}.Informational()

	for _, doc := range docs {
		switch swErr := outbound.rsClient.Index([]redisearch.Document{*doc}...); {
		case mod_errors.Contains(swErr, mod_errors.EDocExist):
			l.Z{l.E: swErr, l.M: doc.Properties[string(_dn)]}.Debug()
		case swErr != nil:
			return swErr
		}
	}

	switch a, b, c := outbound.getDocsByKV(_dn, "dc=domain,dc=tld"); {
	case c != nil:
		return c
	case b != 1 || len(a) != 1:
		fmt.Printf("%v\n", b)
	default:
		fmt.Printf("%v\n", a)
	}

	switch a, b, c := outbound.getDocsByKV(_cn, ""); {
	case c != nil:
		return c
	case b != 1 || len(a) != 1:
		fmt.Printf("%v\n", b)
	default:
		fmt.Printf("%v\n", a)
	}

	return
}
