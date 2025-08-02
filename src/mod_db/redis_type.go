package mod_db

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
)

type _MFV []_FV

type _FV struct {
	_F entryFieldName
	_V string
}

// RedisRepository provides methods for interacting with Redis using rueidis.
type RedisRepository struct {
	client rueidis.Client
	entry  om.Repository[Entry]
	cert   om.Repository[Certificate]
}
