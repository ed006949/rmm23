package mod_db

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
)

// NewRedisRepository creates a new RedisRepository.
func NewRedisRepository(client rueidis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
		entry:  om.NewJSONRepository[Entry](entryKeyHeader, Entry{}, client, om.WithIndexName(entryKeyHeader)),
		cert:   om.NewJSONRepository[Cert](certKeyHeader, Cert{}, client, om.WithIndexName(certKeyHeader)),
		issued: om.NewJSONRepository[Cert](certKeyHeader, Cert{}, client, om.WithIndexName(certKeyHeader)),
	}
}

func searchQueryCommand(search om.FtSearchIndex, query string) rueidis.Completed {
	return search.Query(query).
		Limit().OffsetNum(0, connMaxPaging).
		Build()
}
