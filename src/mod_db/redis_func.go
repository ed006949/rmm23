package mod_db

import (
	"fmt"
	"reflect"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/vmihailenco/msgpack/v5" // MsgPack library for payload

	"rmm23/src/mod_errors"
	"rmm23/src/mod_reflect"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
)

func buildRedisearchSchema(inbound interface{}) (outbound *redisearch.Schema, err error) {
	var (
		schema = redisearch.NewSchema(redisearch.DefaultOptions)
		rt     reflect.Type
	)

	switch rt, err = mod_reflect.GetStructRT(inbound); {
	case err != nil:
		return
	}

	for i := 0; i < rt.NumField(); i++ {
		var (
			field    = rt.Field(i)
			redisTag = field.Tag.Get(redisTagName)
		)

		switch {
		case len(redisTag) == 0:
			continue
		}

		var (
			redisearchTag = field.Tag.Get(rediSearchTagName)
		)

		switch {
		case len(redisearchTag) == 0:
			continue
		}

		var (
			parts = mod_slices.SplitString(redisearchTag, mod_strings.TagSeparator, mod_slices.FlagNormalize)
		)

		switch {
		case len(parts) == 0:
			continue
		}

		var (
			types    = make(map[string]bool)
			options  = make(map[string]bool)
			unknowns = make(map[string]bool)
		)

		for _, opt := range parts {
			switch opt {
			case rediSearchTagTypeIgnore, rediSearchTagTypeText, rediSearchTagTypeNumeric, rediSearchTagTypeTag, rediSearchTagTypeGeo:
				types[opt] = true
			case rediSearchTagOptionSortable:
				options[opt] = true
			default:
				unknowns[opt] = true
			}
		}

		switch {
		case len(types) > 1:
			return nil, mod_errors.ETagMultiType
		case len(unknowns) > 0:
			return nil, mod_errors.ETagUnknown
		}

		switch {
		case types[rediSearchTagTypeIgnore]:
		case types[rediSearchTagTypeText]:
			schema.AddField(redisearch.NewTextFieldOptions(redisTag, redisearch.TextFieldOptions{
				Sortable: options[rediSearchTagOptionSortable],
			}))
		case types[rediSearchTagTypeNumeric]:
			schema.AddField(redisearch.NewNumericFieldOptions(redisTag, redisearch.NumericFieldOptions{
				Sortable: options[rediSearchTagOptionSortable],
			}))
		case types[rediSearchTagTypeTag]:
			schema.AddField(redisearch.NewTagFieldOptions(redisTag, redisearch.TagFieldOptions{
				Sortable:  options[rediSearchTagOptionSortable],
				Separator: mod_strings.SliceSeparator[0],
			}))
		case types[rediSearchTagTypeGeo]:
			schema.AddField(redisearch.NewGeoFieldOptions(redisTag, redisearch.GeoFieldOptions{}))
		default:
			return nil, mod_errors.EUnwilling
		}
	}

	return schema, nil
}

func newRedisearchDocument(schema *redisearch.Schema, docID string, score float32, data interface{}, includePayload bool) (outbound *redisearch.Document, err error) {
	var (
		rv reflect.Value
	)

	switch rv, err = mod_reflect.GetStructRV(data); {
	case err != nil:
		return
	}

	var (
		doc = redisearch.NewDocument(docID, score)
	)

	for i := 0; i < rv.NumField(); i++ {
		var (
			structField = rv.Field(i)
			typeField   = rv.Type().Field(i)
			redisTag    = typeField.Tag.Get(redisTagName)
		)

		switch {
		case len(redisTag) == 0:
			continue
		}

		var (
			schemaField *redisearch.Field
		)

		for _, sf := range schema.Fields {
			switch sf.Name {
			case redisTag:
				schemaField = &sf

				break
			}
		}

		switch {
		case schemaField == nil || !structField.CanInterface():
			continue
		case structField.Kind() == reflect.Ptr && structField.IsNil():
			continue
		case structField.IsZero():
			continue
		}

		var (
			fieldValue interface{}
		)

		switch schemaField.Type {
		case redisearch.TagField, redisearch.TextField:
			fieldValue = fmt.Sprintf("%v", structField.Interface())
		case redisearch.NumericField:
			switch structField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue = float64(structField.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue = float64(structField.Uint())
			case reflect.Float32, reflect.Float64:
				fieldValue = structField.Float()
			case reflect.Bool:
				switch {
				case structField.Bool():
					fieldValue = 1.0
				default:
					fieldValue = 0.0
				}
			default:
				fieldValue = fmt.Sprintf("%v", structField.Interface())
			}
		default:
			fieldValue = fmt.Sprintf("%v", structField.Interface())
		}

		doc.Set(schemaField.Name, fieldValue)
	}

	switch {
	case includePayload:
		var (
			encodedPayload []byte
		)

		switch encodedPayload, err = msgpack.Marshal(data); {
		case err != nil:
			return nil, err
		}

		doc.SetPayload(encodedPayload)
	}

	return &doc, nil
}
