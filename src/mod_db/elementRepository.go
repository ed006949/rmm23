package mod_db

import (
	"context"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
)

// RedisRepository provides methods for interacting with Redis using rueidis.
type RedisRepository struct {
	ctx    context.Context
	client rueidis.Client
	info   map[string]*ftInfo
	entry  om.Repository[Entry]
	cert   om.Repository[Cert]
	// issued    om.Repository[Cert]
}
