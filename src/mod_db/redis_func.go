package mod_db

import (
	"strings"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
)

// RedisRepository provides methods for interacting with Redis using rueidis.
type RedisRepository struct {
	repo om.Repository[Entry]
}

// NewRedisRepository creates a new RedisRepository.
func NewRedisRepository(client rueidis.Client) *RedisRepository {
	return &RedisRepository{
		repo: om.NewJSONRepository[Entry](entryKeyHeader, Entry{}, client, om.WithIndexName(entryKeyHeader)),
	}
}

func escapeQueryValue(inbound string) string {
	replacer := strings.NewReplacer(
		"=", "\\=",
		",", "\\,",
		// "(", "\\(",
		// ")", "\\)",
		// "{", "\\{",
		// "}", "\\}",
		// "[", "\\[",
		// "]", "\\]",
		"\"", "\\\"",
		"'", "\\'",
		"~", "\\~",
	)

	return replacer.Replace(inbound)
}
