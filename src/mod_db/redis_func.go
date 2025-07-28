package mod_db

import (
	"context"
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
		repo: om.NewJSONRepository[Entry](entryKeyHeader, Entry{}, client),
	}
}

// GetEntry retrieves an Entry from Redis.
func GetEntry(ctx context.Context, repo *RedisRepository, id string) (*Entry, error) {
	return repo.FindEntry(ctx, id)
}

// SetEntry saves an Entry to Redis.
func SetEntry(ctx context.Context, repo *RedisRepository, e *Entry) error {
	return repo.SaveEntry(ctx, e)
}

// DelEntry saves an Entry to Redis.
func DelEntry(ctx context.Context, repo *RedisRepository, id string) error {
	return repo.repo.Remove(ctx, id)
}

func escapeQueryValue(inbound string) string {
	replacer := strings.NewReplacer(
		"=", "\\=",
		",", "\\,",
		"(", "\\(",
		")", "\\)",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"\"", "\\\"",
		"'", "\\'",
		"~", "\\~",
	)

	return replacer.Replace(inbound)
}
