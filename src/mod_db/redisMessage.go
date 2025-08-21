package mod_db

import (
	"github.com/redis/rueidis"

	"rmm23/src/mod_errors"
)

func parseRedisMessages(messages map[string]rueidis.RedisMessage) (outbound map[string]any, err error) {
	outbound = make(map[string]any, len(messages))

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
	case message.IsArray():
		return parseRedisMessageArray(message)

	case message.IsBool():
		return message.ToBool()

	case message.IsFloat64():
		return message.ToFloat64()

	case message.IsInt64():
		return message.ToInt64()

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

func parseRedisMessageArray(message rueidis.RedisMessage) (outbound any, err error) {
	// Try as map first (for Redis hash responses)
	switch messages, mapErr := message.AsMap(); {
	case mapErr == nil:
		return parseRedisMessages(messages)
	}

	// Parse as regular array
	switch messages, arrErr := message.ToArray(); {
	case arrErr == nil:
		var (
			interim = make([]any, 0, len(messages))
		)

		for _, b := range messages {
			switch message2, forErr := parseRedisMessage(b); {
			case forErr != nil:
				return nil, forErr
			default:
				interim = append(interim, message2)
			}
		}

		return interim, nil
	}

	return nil, mod_errors.EParse
}
