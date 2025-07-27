package mod_db

import (
	"fmt"
	"reflect"

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

	// Handle pointer types
	switch targetType.Kind() {
	case reflect.Ptr:
		switch {
		case fieldValue.IsNil():
			fieldValue.Set(reflect.New(targetType.Elem()))
		}

		fieldValue = fieldValue.Elem()
		targetType = targetType.Elem()
	}

	// Check for RedisUnmarshaler interface
	switch unmarshaler, ok := fieldValue.Addr().Interface().(RedisUnmarshaler); {
	case ok:
		return unmarshaler.UnmarshalRedis(propValue)
	default:
		return fmt.Errorf("field type %s does not implement RedisUnmarshaler interface", targetType.Kind())
	}
}
