package mod_db

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
)

type MFV []FV

type FV struct {
	Field entryFieldName
	Value string
}

// RedisRepository provides methods for interacting with Redis using rueidis.
type RedisRepository struct {
	client rueidis.Client
	entry  om.Repository[Entry]
	cert   om.Repository[Cert]
}
