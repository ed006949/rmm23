package mod_db

import (
	"github.com/redis/rueidis"

	"rmm23/src/mod_net"
)

type Conf struct {
	URL    *mod_net.URL `json:"url,omitempty"`
	Name   string       `json:"name,omitempty"`
	client rueidis.Client
	repo   *RedisRepository
	// rcNetwork string
	// rsClient        *redisearch.Client
	// schema          *redisearch.Schema
	// schemaMap       schemaMapType
	// indexDefinition *redisearch.IndexDefinition
}

type entryFieldName string

// type schemaMapType map[entryFieldName]redisearch.FieldType
