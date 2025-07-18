package mod_db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/vmihailenco/msgpack/v5" // MsgPack library for payload

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

func newDocumentFromStruct(schema *redisearch.Schema, docID string, score float32, data interface{}, includePayload bool) (redisearch.Document, error) {
	doc := redisearch.NewDocument(docID, score)
	val := reflect.ValueOf(data)

	// Ensure data is a pointer to a struct, and get the underlying struct value
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return redisearch.Document{}, fmt.Errorf("data must be a struct or a pointer to a struct, got %T", data)
	}

	// Iterate over struct fields to find corresponding schema fields
	for i := 0; i < val.NumField(); i++ {
		structField := val.Field(i)
		typeField := val.Type().Field(i)

		redisTag := typeField.Tag.Get(redisTagName)
		if redisTag == "" {
			continue // Skip fields without a redis tag
		}

		// Construct the schema field name as it would be in the schema
		schemaFieldName := "$." + redisTag

		// Find the corresponding schema field
		var schemaField *redisearch.Field
		for _, sf := range schema.Fields {
			if sf.Name == schemaFieldName {
				schemaField = sf
				break
			}
		}

		if schemaField == nil || !structField.CanInterface() {
			// Schema field not found or struct field not exportable, skip
			continue
		}

		var fieldValue interface{}
		switch schemaField.Type {
		case redisearch.TagField, redisearch.TextField:
			// For text and tag fields, convert any type to string
			fieldValue = fmt.Sprintf("%v", structField.Interface())
		case redisearch.NumericField:
			// Handle numeric types (int, float, bool to float64)
			switch structField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue = float64(structField.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue = float64(structField.Uint())
			case reflect.Float32, reflect.Float64:
				fieldValue = structField.Float()
			case reflect.Bool: // Convert bool to 0.0 or 1.0 for numeric
				if structField.Bool() {
					fieldValue = 1.0
				} else {
					fieldValue = 0.0
				}
			default:
				// Fallback for other numeric-like types, try to convert to string
				// RediSearch might handle string-to-numeric conversion if valid
				fieldValue = fmt.Sprintf("%v", structField.Interface())
			}
		// Add cases for other RediSearch field types (Geo, etc.) if needed
		default:
			// Default to string conversion for unsupported types
			fieldValue = fmt.Sprintf("%v", structField.Interface())
		}

		doc.Set(schemaField.Name, fieldValue)
	}

	// Handle payload if requested
	if includePayload {
		encodedPayload, err := msgpack.Marshal(data)
		if err != nil {
			return redisearch.Document{}, fmt.Errorf("failed to marshal payload: %w", err)
		}
		doc.SetPayload(encodedPayload)
	}

	return doc, nil
}
