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
	info   map[string]*FTInfo
	entry  om.Repository[Entry]
	cert   om.Repository[Cert]
	// issued    om.Repository[Cert]
}

// type REntries om.Repository[Entry]
// type RCerts om.Repository[Cert]
// type RIssued om.Repository[Cert]
