package mod_db

import (
	"github.com/RediSearch/redisearch-go/redisearch"

	"rmm23/src/mod_net"
)

type Conf struct {
	URL             *mod_net.URL `json:"url,omitempty"`
	Name            string       `json:"name,omitempty"`
	rcNetwork       string
	rsClient        *redisearch.Client
	schema          *redisearch.Schema
	schemaMap       schemaMapType
	indexDefinition *redisearch.IndexDefinition
}

type entryFieldName string

type schemaMapType map[entryFieldName]redisearch.FieldType
