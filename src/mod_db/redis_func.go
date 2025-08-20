package mod_db

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/mod_errors"
)

// NewRedisRepository creates a new RedisRepository.
func NewRedisRepository(client rueidis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
		entry:  om.NewJSONRepository[Entry](entryKeyHeader, Entry{}, client, om.WithIndexName(entryKeyHeader)),
		cert:   om.NewJSONRepository[Cert](certKeyHeader, Cert{}, client, om.WithIndexName(certKeyHeader)),
	}
}

func searchQueryCommand(search om.FtSearchIndex, query string) rueidis.Completed {
	return search.Query(query).
		Limit().OffsetNum(0, connMaxPaging).
		Build()
}

func parseRedisMessages(messages map[string]rueidis.RedisMessage) (outbound map[string]any, err error) {
	outbound = make(map[string]any)

	for c, message := range messages {
		switch outbound[c], err = parseRedisMessage(message); {
		case err != nil:
			return nil, err
		}
	}

	return
}

func parseRedisMessage(message rueidis.RedisMessage) (outbound any, err error) {
	switch {
	// case message.IsCacheHit():
	case message.IsArray():
		switch messages, swErr := message.AsMap(); {
		case swErr == nil:
			return parseRedisMessages(messages)
		}

		switch messages, swErr := message.ToArray(); {
		case swErr == nil:
			var (
				interim []any
			)

			for _, b := range messages {
				switch message2, swErr2 := parseRedisMessage(b); {
				case swErr2 != nil:
					return nil, swErr2
				default:
					interim = append(interim, message2)
				}
			}

			return interim, nil
		}

		return nil, mod_errors.EParse

	case message.IsBool():
		return message.AsBool()

	case message.IsFloat64():
		return message.AsFloat64()

	case message.IsInt64():
		return message.AsInt64()

	case message.IsNil():
		return nil, nil

	case message.IsMap():
		return message.ToMap()

	case message.IsString():
		return message.ToString()

	default:
		return message.ToAny()
	}
}
