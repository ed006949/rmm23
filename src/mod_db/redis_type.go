package mod_db

import (
	"github.com/redis/rueidis/om"
)

type _FVOF struct {
	_FV []_FV
	_OF []entryFieldName
}

type _FV struct {
	_F entryFieldName
	_V string
}

// RedisRepository provides methods for interacting with Redis using rueidis.
type RedisRepository struct {
	repo om.Repository[Entry]
}
