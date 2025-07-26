package mod_db

import (
	"context"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
)

// func (e *ElementDomain) redisearchSchema() *redisearch.Schema { return buildRedisearchSchema(e) }
// func (e *ElementGroup) redisearchSchema() *redisearch.Schema  { return buildRedisearchSchema(e) }
// func (e *ElementUser) redisearchSchema() *redisearch.Schema   { return buildRedisearchSchema(e) }
// func (e *ElementHost) redisearchSchema() *redisearch.Schema   { return buildRedisearchSchema(e) }

func (r *Entry) redisearchSchema() (outbound *redisearch.Schema, err error) {
	return buildRedisearchSchema(r)
}

func (r *Conf) dial() (err error) {
	switch r.rcNetwork, err = r.URL.RedisNetwork(); {
	case err != nil:
		return
	}

	r.rsClient = redisearch.NewClientFromPool(&redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialContext(ctx, r.rcNetwork, r.URL.Host, redis.DialDatabase(0))
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
	}, r.Name)

	switch {
	case r.rsClient == nil:
		return mod_errors.ENoConn
	}

	// _ = r.rsClient.Drop() // test&dev, delete everything
	_ = r.rsClient.DropIndex(false) // prod, delete index only

	return
}

func (r *Conf) getDoc(inbound string) (outbound *redisearch.Document, err error) {
	return r.rsClient.Get(inbound)
}

func (r *Conf) getDocByUUID(inbound mod_ldap.AttrUUID) (outbound *redisearch.Document, err error) {
	return r.getDoc(inbound.Entry())
}

func (r *Conf) getDocsByKV(key entryFieldName, value any) (outbound []redisearch.Document, count int, err error) {
	return r.rsClient.Search(redisearch.NewQuery(createQuery(key, value)).Limit(0, connMaxPaging))
}
