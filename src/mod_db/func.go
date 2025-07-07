package mod_db

import (
	"reflect"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
)

// buildRedisearchSchema dynamically builds a RediSearch schema from a struct's tags.
// It uses the `redis` tag for the field name and the `redisearch` tag for the type and options.
// Example: `redis:"user_name" redisearch:"text,sortable"`
func buildRedisearchSchema(s interface{}) *redisearch.Schema {
	var (
		schema = redisearch.NewSchema(redisearch.DefaultOptions)
		elem   = reflect.TypeOf(s).Elem()
	)

	for i := 0; i < elem.NumField(); i++ {
		var (
			field    = elem.Field(i)
			redisTag = field.Tag.Get("redis")
		)
		switch redisTag {
		case "":
			continue
		}

		var (
			redisearchTag = field.Tag.Get("redisearch")
		)
		switch redisearchTag {
		case "":
			continue
		}

		// Parse the redisearch tag string.
		var (
			parts = strings.Split(redisearchTag, ",")
		)
		switch len(parts) {
		case 0:
			continue
		}

		var (
			types    = make(map[string]bool)
			options  = make(map[string]bool)
			unknowns = make(map[string]bool)
		)
		for _, opt := range parts {
			var (
				trimmedOpt = strings.TrimSpace(opt)
			)
			switch trimmedOpt {
			case "-", "text", "numeric":
				types[trimmedOpt] = true
			case "sortable":
				options[trimmedOpt] = true
			default:
				unknowns[trimmedOpt] = true
			}
		}

		// check for tag-errors
		switch {
		case len(types) > 1:
			panic("multiple types")
		case len(unknowns) > 0:
			panic("unknown options")
		}

		// parsing types
		switch {
		case types["-"]:
		case types["text"]:
			schema.AddField(redisearch.NewTextFieldOptions(redisTag, redisearch.TextFieldOptions{
				Sortable: options["sortable"],
			}))
		case types["numeric"]:
			schema.AddField(redisearch.NewNumericFieldOptions(redisTag, redisearch.NumericFieldOptions{
				Sortable: options["sortable"],
			}))
		}
	}
	return schema
}
