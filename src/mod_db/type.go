package mod_db

import (
	"rmm23/src/mod_net"
)

type Conf struct {
	URL  *mod_net.URL `json:"url,omitempty"`
	repo *RedisRepository
}

type entryFieldName string

// type schemaMapType map[entryFieldName]redisearch.FieldType
