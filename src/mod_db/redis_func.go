package mod_db

import (
	"reflect"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
)

// buildRedisearchSchema dynamically builds a RediSearch schema from a struct's tags.
// It uses the `redis` tag for the field name and the `redisearch` tag for the type and options.
// Example: `redis:"user_name" redisearch:"text,sortable"`
func buildRedisearchSchema(inbound interface{}) *redisearch.Schema {
	var (
		schema = redisearch.NewSchema(redisearch.DefaultOptions) // Initialize a new schema with default options.
		elem   = reflect.TypeOf(inbound).Elem()                  // Get the type of the element that the interface `inbound` points to.
	)

	// Iterate over all the fields of the struct.
	for i := 0; i < elem.NumField(); i++ {
		var (
			field    = elem.Field(i)
			redisTag = field.Tag.Get(redisTagName) // Get the `redis` tag, which defines the field name in Redis.
		)
		switch redisTag {
		case "": // Skip fields that don't have a `redis` tag.
			continue
		}

		var (
			redisearchTag = field.Tag.Get(rediSearchTagName) // Get the `redisearch` tag, which defines the field type and options for RediSearch.
		)
		switch redisearchTag {
		case "": // Skip fields that don't have a `redisearch` tag.
			continue
		}

		// Parse the redisearch tag string, which is expected to be comma-separated.
		var (
			parts = strings.Split(redisearchTag, ",")
		)
		// If the tag is empty after splitting, skip this field.
		switch len(parts) {
		case 0:
			continue
		}

		// maps to hold the parsed tag components.
		var (
			types    = make(map[string]bool) // e.g., "text", "numeric"
			options  = make(map[string]bool) // e.g., "sortable"
			unknowns = make(map[string]bool) // for any unrecognized parts
		)
		// Categorize each part of the tag.
		for _, opt := range parts {
			var (
				trimmedOpt = strings.TrimSpace(opt)
			)
			switch trimmedOpt {
			case rediSearchTagTypeIgnore, rediSearchTagTypeText, rediSearchTagTypeNumeric:
				types[trimmedOpt] = true
			case rediSearchTagOptionSortable:
				options[trimmedOpt] = true
			default:
				unknowns[trimmedOpt] = true
			}
		}

		// Check for errors in the tag definition.
		switch {
		case len(types) > 1: // A field can only have one type.
			panic("multiple types")
		case len(unknowns) > 0: // The tag should not contain any unknown options.
			panic("unknown tag fields")
		}

		// Add the field to the schema based on its type.
		switch {
		case types[rediSearchTagTypeIgnore]: // The "-" type indicates that the field should be ignored.
		case types[rediSearchTagTypeText], types[rediSearchTagTypeNumeric]: // Handle 'text', 'numeric' type fields.
			schema.AddField(redisearch.NewTagFieldOptions("$."+redisTag, redisearch.TagFieldOptions{
				Sortable:  options[rediSearchTagOptionSortable],
				Separator: '\x1f',
			}))
			// case types[rediSearchTagTypeNumeric]: // Handle 'numeric' type fields.
			// 	schema.AddField(redisearch.NewNumericFieldOptions("$."+redisTag, redisearch.NumericFieldOptions{
			// 		Sortable: options[rediSearchTagOptionSortable],
			// 	}))
		}
	}

	// Return the fully constructed schema.
	return schema
}
