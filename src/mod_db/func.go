package mod_db

import (
	"context"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_slices"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.LDAPConfig, outbound *Conf) (err error) {
	// switch err = inbound.Search(); {
	// case err != nil:
	// 	return
	// }
	var (
		docs   []*redisearch.Document
		schema = new(Entry).redisearchSchema()
	)

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

	var (
		// RediSearch requires DB 0 for index creation
		// rcDB      = 0

		rcNetwork string
		rcName    = "entryIdx"
	)

	switch rcNetwork, err = outbound.URL.RedisNetwork(); {
	case err != nil:
		return
	}

	var (
		rsClient = redisearch.NewClientFromPool(&redis.Pool{
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialContext(ctx, rcNetwork, outbound.URL.Host, redis.DialDatabase(0))
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) (tErr error) {
				_, tErr = c.Do(_PING)

				return
			},
			MaxIdle:         connMaxIdle,
			MaxActive:       connMaxActive,
			IdleTimeout:     connIdleTimeout,
			Wait:            connWait,
			MaxConnLifetime: connMaxConnLifetime,
		}, rcName)
		// entry  = Entry{}
		// schema = entry.redisearchSchema()
		// _ = rsClient.Drop() // test&dev, delete old entries
		_ = rsClient.DropIndex(false) // prod, don't delete old entries

		rsQuery = redisearch.NewQuery("*").SetReturnFields("uuid", "dn").Limit(0, 1000000)
	)

	switch err = rsClient.CreateIndex(schema); {
	// case mod_errors.Contains(err, mod_errors.EIndexExist):
	case err != nil:
		return
	}

	switch a, b, c := rsClient.Search(rsQuery); {
	case c != nil:
		return c
	default:
		a = a
		b = b

		panic(nil)
	}

	for _, doc := range docs {
		switch err = rsClient.Index([]redisearch.Document{*doc}...); {
		case mod_errors.Contains(err, mod_errors.EDocExist):
		case err != nil:
			return
		}
	}

	return
}
