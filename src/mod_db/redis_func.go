package mod_db

import (
	"reflect"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"

	"rmm23/src/mod_strings"
)

func buildRedisearchSchema(inbound interface{}) *redisearch.Schema {
	var (
		schema = redisearch.NewSchema(redisearch.DefaultOptions)
		elem   = reflect.TypeOf(inbound).Elem()
	)

	for i := 0; i < elem.NumField(); i++ {
		var (
			field    = elem.Field(i)
			redisTag = field.Tag.Get(redisTagName)
		)
		switch redisTag {
		case "":
			continue
		}

		var (
			redisearchTag = field.Tag.Get(rediSearchTagName)
		)
		switch redisearchTag {
		case "":
			continue
		}

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
			case rediSearchTagTypeIgnore, rediSearchTagTypeText, rediSearchTagTypeNumeric, rediSearchTagTypeTag, rediSearchTagTypeGeo:
				types[trimmedOpt] = true
			case rediSearchTagOptionSortable:
				options[trimmedOpt] = true
			default:
				unknowns[trimmedOpt] = true
			}
		}

		switch {
		case len(types) > 1:
			panic("multiple types")
		case len(unknowns) > 0:
			panic("unknown tag fields")
		}

		switch {
		case types[rediSearchTagTypeIgnore]:
		case types[rediSearchTagTypeText]:
			schema.AddField(redisearch.NewTextFieldOptions("$."+redisTag, redisearch.TextFieldOptions{
				Sortable: options[rediSearchTagOptionSortable],
			}))
		case types[rediSearchTagTypeNumeric]:
			schema.AddField(redisearch.NewNumericFieldOptions("$."+redisTag, redisearch.NumericFieldOptions{
				Sortable: options[rediSearchTagOptionSortable],
			}))
		case types[rediSearchTagTypeTag]:
			schema.AddField(redisearch.NewTagFieldOptions("$."+redisTag, redisearch.TagFieldOptions{
				Sortable:  options[rediSearchTagOptionSortable],
				Separator: mod_strings.SliceDelimiter[0],
			}))
		case types[rediSearchTagTypeGeo]:
			schema.AddField(redisearch.NewGeoFieldOptions("$."+redisTag, redisearch.GeoFieldOptions{}))
		default:
			panic("unwilling to perform")
		}
	}

	return schema
}
