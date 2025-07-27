package mod_db

import (
	"context"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
)

// func (e *ElementDomain) redisearchSchema() *redisearch.Schema { return buildRedisearchSchema(e) }
// func (e *ElementGroup) redisearchSchema() *redisearch.Schema  { return buildRedisearchSchema(e) }
// func (e *ElementUser) redisearchSchema() *redisearch.Schema   { return buildRedisearchSchema(e) }
// func (e *ElementHost) redisearchSchema() *redisearch.Schema   { return buildRedisearchSchema(e) }

func (r *entry) redisearchSchema() (schema *redisearch.Schema, schemaMap schemaMapType, err error) {
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

	return
}

func (r *Conf) createIndex() (err error) {
	var (
		indexInfo *redisearch.IndexInfo
	)

	// define indexDefinition
	r.indexDefinition = redisearch.NewIndexDefinition().AddPrefix(entryDocIDHeader)

	switch l.CLEAR {
	case true:
		_ = r.rsClient.Drop()
	// _ = r.rsClient.DropIndex(true)
	default:
		_ = r.rsClient.DropIndex(false)
	}

	switch swErr := r.rsClient.CreateIndexWithIndexDefinition(r.schema, r.indexDefinition); {
	case mod_errors.Contains(swErr, mod_errors.EIndexExist):
	case swErr != nil:
		return swErr
	}

	// wait for index to complete
	switch indexInfo, err = r.rsClient.Info(); {
	case err != nil:
		return
	}

	for indexInfo.IsIndexing {
		switch indexInfo, err = r.rsClient.Info(); {
		case err != nil:
			return
		}
	}

	return
}

func (r *Conf) getDoc(inbound string) (outbound *redisearch.Document, err error) {
	return r.rsClient.Get(inbound)
}

func (r *Conf) getDocByUUID(inbound attrUUID) (outbound *redisearch.Document, err error) {
	return r.getDoc(inbound.Entry())
}

func (r *Conf) getDocsByKV(key entryFieldName, value string) (outbound []redisearch.Document, count int, err error) {
	var (
		interim string
	)

	switch r.schemaMap[key] {
	case redisearch.TextField:
		interim = "@" + key.String() + ":" + escapeQueryValue(value) + ""
	case redisearch.NumericField:
		interim = "@" + key.String() + ":[" + escapeQueryValue(value) + "]"
	case redisearch.GeoField:
		interim = "@" + key.String() + ":[" + escapeQueryValue(value) + "]"
	case redisearch.TagField:
		interim = "@" + key.String() + ":{" + escapeQueryValue(value) + "}"
	default:
		return nil, 0, mod_errors.EUnwilling
	}

	return r.rsClient.Search(redisearch.NewQuery(interim).SetInFields(key.String()).Limit(0, connMaxPaging))
}
