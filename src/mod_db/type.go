package mod_db

import (
	"rmm23/src/mod_url"
)

type Conf struct {
	URL  *mod_url.URL     `json:"url,omitempty"`
	Repo *RedisRepository `json:"-"`
}

// type schemaMapType map[EntryFieldName]redisearch.FieldType
