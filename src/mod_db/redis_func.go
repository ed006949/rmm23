package mod_db

import (
	"strings"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
)

func escapeQueryValue(inbound string) (outbound string) {
	replacer := strings.NewReplacer(
		"=", "\\=", //
		",", "\\,", //
		"(", "\\(", //
		")", "\\)", //
		"{", "\\{", //
		"}", "\\}", //
		"[", "\\[", //
		"]", "\\]", //
		"\"", "\\\"", //
		"'", "\\'", //
		"~", "\\~", //
		"-", "\\-", // (?)
	)

	return replacer.Replace(inbound)
}

// NewRedisRepository creates a new RedisRepository.
func NewRedisRepository(client rueidis.Client) *RedisRepository {
	return &RedisRepository{
		repo: om.NewJSONRepository[Entry](entryKeyHeader, Entry{}, client, om.WithIndexName(entryKeyHeader)),
	}
}
