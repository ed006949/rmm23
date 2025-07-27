package mod_db

import (
	"fmt"
	"reflect"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/vmihailenco/msgpack/v5"

	"rmm23/src/mod_reflect"
	"rmm23/src/mod_slices"
)

func marshalRedisearchDoc(schema *redisearch.Schema, docID string, score float32, data interface{}, includePayload bool) (outbound *redisearch.Document, err error) {
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
			redisTag = rv.Type().Field(i).Tag.Get(redisTagName)
		)

		switch {
		case len(redisTag) == 0:
			continue
		}

		var (
			structField = rv.Field(i)
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

		doc.Set(schemaField.Name, getFieldValue(schemaField, &structField))
	}

	switch err = setPayload(&doc, data, includePayload); {
	case err != nil:
		return
	}

	return &doc, nil
}

func setPayload(doc *redisearch.Document, data interface{}, includePayload bool) (err error) {
	switch {
	case includePayload:
		var (
			encodedPayload []byte
		)

		switch encodedPayload, err = msgpack.Marshal(data); {
		case err != nil:
			return
		}

		doc.SetPayload(encodedPayload)
	}

	return
}

func getFieldValue(schemaField *redisearch.Field, structField *reflect.Value) (fieldValue any) {
	switch schemaField.Type {
	case redisearch.TagField:
		switch structField.Kind() {
		case reflect.Slice, reflect.Array:
			var (
				elements []string
			)
			for i := 0; i < structField.Len(); i++ {
				elements = append(elements, fmt.Sprintf("%v", structField.Index(i).Interface()))
			}

			fieldValue = mod_slices.Join(elements, string(sliceSeparator), mod_slices.FlagNormalize)
			// fieldValue = strings.Join(elements, string(sliceSeparator))
		default:
			fieldValue = fmt.Sprintf("%v", structField.Interface())
		}
	case redisearch.TextField:
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

	return
}
