package mod_db

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
)

// RedisRepository provides methods for interacting with Redis using rueidis.
type RedisRepository struct {
	client rueidis.Client
	entry  om.Repository[Entry]
	cert   om.Repository[Cert]
	issued om.Repository[Cert]
}
