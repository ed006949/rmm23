package mod_db

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/vmihailenco/msgpack/v5"
)

func unmarshalRedisearchDoc(doc *redisearch.Document, outbound interface{}) (err error) {
	switch {
	case doc == nil:
		return fmt.Errorf("input redisearch.Document is nil")
	case reflect.TypeOf(outbound).Kind() != reflect.Ptr || reflect.TypeOf(outbound).Elem().Kind() != reflect.Struct:
		return fmt.Errorf("outbound must be a pointer to a struct, got %T", outbound)
	}

	switch {
	case len(doc.Payload) > 0:
		switch err = msgpack.Unmarshal(doc.Payload, outbound); {
		case err != nil:
			return fmt.Errorf("failed to unmarshal payload for document ID '%s': %w", doc.Id, err)
		}
	}

	var (
		sv = reflect.ValueOf(outbound).Elem()
		st = sv.Type()
	)

	for i := 0; i < st.NumField(); i++ {
		var (
			ft = st.Field(i)
			fv = sv.Field(i)
		)

		switch {
		case len(ft.PkgPath) != 0:
			continue
		}

		var (
			redisTag = ft.Tag.Get(redisTagName)
		)
		switch {
		case len(redisTag) == 0:
			continue
		}

		var (
			propValue, ok = doc.Properties[redisTag]
		)
		switch {
		case !ok:
			continue
		}

		switch err = setStructFieldValue(fv, propValue, ft.Type); {
		case err != nil:
			return fmt.Errorf("failed to set field '%s' (Redis tag '%s') for document ID '%s': %w", ft.Name, redisTag, doc.Id, err)
		}
	}

	return nil
}

func setStructFieldValue(fieldValue reflect.Value, propValue interface{}, targetType reflect.Type) (err error) {
	switch {
	case !fieldValue.CanSet():
		return fmt.Errorf("cannot set unexported field '%s'", fieldValue.Type().Name())
	}

	switch targetType.Kind() {
	case reflect.Ptr:
		switch {
		case fieldValue.IsNil():
			fieldValue.Set(reflect.New(targetType.Elem()))
		}

		fieldValue = fieldValue.Elem()
		targetType = targetType.Elem()
	}

	switch targetType.Kind() {
	case reflect.String:
		switch s, ok := propValue.(string); {
		case ok:
			fieldValue.SetString(s)
		default:
			return fmt.Errorf("expected string for field type '%s', got %T", targetType.Kind(), propValue)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch f, ok := propValue.(float64); {
		case ok:
			fieldValue.SetInt(int64(f))
		default:
			return fmt.Errorf("expected numeric (float64) for field type '%s', got %T", targetType.Kind(), propValue)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch f, ok := propValue.(float64); {
		case ok:
			fieldValue.SetUint(uint64(f))
		default:
			return fmt.Errorf("expected numeric (float64) for field type '%s', got %T", targetType.Kind(), propValue)
		}

	case reflect.Float32, reflect.Float64:
		switch f, ok := propValue.(float64); {
		case ok:
			fieldValue.SetFloat(f)
		default:
			return fmt.Errorf("expected float64 for field type '%s', got %T", targetType.Kind(), propValue)
		}

	case reflect.Bool:
		switch {
		case propValue == nil:
			return fmt.Errorf("expected bool, string, or numeric for field type '%s', got nil", targetType.Kind())
		case reflect.TypeOf(propValue).Kind() == reflect.Bool:
			fieldValue.SetBool(propValue.(bool))
		case reflect.TypeOf(propValue).Kind() == reflect.String:
			var (
				parsedBool bool
			)
			switch parsedBool, err = strconv.ParseBool(propValue.(string)); {
			case err != nil:
				return fmt.Errorf("failed to parse bool string '%s': %w", propValue.(string), err)
			default:
				fieldValue.SetBool(parsedBool)
			}
		case reflect.TypeOf(propValue).Kind() == reflect.Float64:
			fieldValue.SetBool(propValue.(float64) != 0.0)
		default:
			return fmt.Errorf("expected bool, string, or numeric for field type '%s', got %T", targetType.Kind(), propValue)
		}

	case reflect.Slice, reflect.Array:
		switch s, ok := propValue.(string); {
		case ok:
			var (
				elements      = strings.Split(s, string(sliceSeparator))
				sliceElemType = targetType.Elem()
				newSlice      = reflect.MakeSlice(targetType, len(elements), len(elements))
			)
			for i, elemStr := range elements {
				var (
					elemVal = reflect.New(sliceElemType).Elem()
				)
				switch err = setStructFieldValue(elemVal, elemStr, sliceElemType); {
				case err != nil:
					return fmt.Errorf("failed to unmarshal slice element %d ('%s') for field type '%s': %w", i, elemStr, targetType.Kind(), err)
				default:
					newSlice.Index(i).Set(elemVal)
				}
			}

			fieldValue.Set(newSlice)
		default:
			return fmt.Errorf("expected string for slice/array field type '%s', got %T", targetType.Kind(), propValue)
		}

	case reflect.Map:
		switch s, ok := propValue.(string); {
		case ok:
			var (
				newMap = reflect.MakeMap(targetType)
			)
			switch err = json.Unmarshal([]byte(s), newMap.Addr().Interface()); {
			case err != nil:
				return fmt.Errorf("failed to unmarshal JSON string to map for field type '%s': %w", targetType.Kind(), err)
			default:
				fieldValue.Set(newMap)
			}
		default:
			return fmt.Errorf("expected string for map field type '%s', got %T", targetType.Kind(), propValue)
		}

	default:
		return fmt.Errorf("unsupported target field type for unmarshaling: %s (field value type %T)", targetType.Kind(), propValue)
	}

	return nil
}
