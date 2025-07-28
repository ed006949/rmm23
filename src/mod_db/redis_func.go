package mod_db

import (
	"context"
	"strings"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
)

// RedisRepository provides methods for interacting with Redis using rueidis.
type RedisRepository struct {
	repo om.Repository[entry]
}

// NewRedisRepository creates a new RedisRepository.
func NewRedisRepository(client rueidis.Client) *RedisRepository {
	return &RedisRepository{
		repo: om.NewJSONRepository[entry](entryPrefix, entry{}, client),
	}
}

// GetEntry retrieves an entry from Redis.
func GetEntry(ctx context.Context, repo *RedisRepository, id string) (*entry, error) {
	return repo.FindEntry(ctx, id)
}

// SetEntry saves an entry to Redis.
func SetEntry(ctx context.Context, repo *RedisRepository, e *entry) error {
	return repo.SaveEntry(ctx, e)
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
