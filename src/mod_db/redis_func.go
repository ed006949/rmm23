package mod_db

import (
	"fmt"
	"strings"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"
)

// NewRedisRepository creates a new RedisRepository.
func NewRedisRepository(client rueidis.Client) *RedisRepository {
	return &RedisRepository{
		entry: om.NewJSONRepository[Entry](entryKeyHeader, Entry{}, client, om.WithIndexName(entryKeyHeader)),
		cert:  om.NewJSONRepository[Certificate](certKeyHeader, Certificate{}, client, om.WithIndexName(certKeyHeader)),
	}
}

func buildFVQuery(field entryFieldName, value string) (outbound string) {
	return fmt.Sprintf(
		"@%s:%s%v%s",
		field.String(),
		entryFieldValueEnclosure[entryFieldMap[field]][0],
		escapeQueryValue(value),
		entryFieldValueEnclosure[entryFieldMap[field]][1],
	)
}

func escapeQueryValue(inbound string) (outbound string) {
	replacer := strings.NewReplacer(
		`=`, `\=`, //
		`,`, `\,`, //
		`(`, `\(`, //
		`)`, `\)`, //
		`{`, `\{`, //
		`}`, `\}`, //
		`[`, `\[`, //
		`]`, `\]`, //
		`"`, `\"`, //
		`'`, `\'`, //
		`~`, `\~`, //
		`-`, `\-`, // (?)
	)

	return replacer.Replace(inbound)
}

func searchQueryCommand(search om.FtSearchIndex, query string) rueidis.Completed {
	return search.Query(query).
		Limit().OffsetNum(0, connMaxPaging).
		Build()
}
