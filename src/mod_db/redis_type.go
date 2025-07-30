package mod_db

import (
	"github.com/redis/rueidis/om"
)

type _FV struct {
	_F entryFieldName
	_V string
}

// RedisRepository provides methods for interacting with Redis using rueidis.
type RedisRepository struct {
	repo om.Repository[Entry]
}
