package mod_db

import (
	"github.com/RediSearch/redisearch-go/redisearch"
)

func (e *ElementDomain) RedisearchSchema() *redisearch.Schema { return buildRedisearchSchema(e) }
func (e *ElementGroup) RedisearchSchema() *redisearch.Schema  { return buildRedisearchSchema(e) }
func (e *ElementUser) RedisearchSchema() *redisearch.Schema   { return buildRedisearchSchema(e) }
func (e *ElementHost) RedisearchSchema() *redisearch.Schema   { return buildRedisearchSchema(e) }
