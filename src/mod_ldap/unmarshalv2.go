package mod_ldap

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// WalkTags processes struct fields with encoding.TextUnmarshaler support and standard type fallbacks.
func WalkTags(entries []*ldap.Entry, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("target must be pointer to slice")
	}

	sliceValue := targetValue.Elem()
	sliceType := sliceValue.Type().Elem()

	for _, entry := range entries {
		item := reflect.New(sliceType).Elem()

		if err := walkStructFields(entry, item); err != nil {
			return fmt.Errorf("failed to process entry %s: %w", entry.DN, err)
		}

		sliceValue.Set(reflect.Append(sliceValue, item))
	}

	return nil
}

func walkStructFields(entry *ldap.Entry, structValue reflect.Value) error {
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		ldapTag := field.Tag.Get("ldap")
		if ldapTag == "" {
			continue
		}

		// Handle DN specially
		if ldapTag == "dn" {
			if err := setFieldValue(fieldValue, []string{entry.DN}); err != nil {
				return fmt.Errorf("failed to set DN field %s: %w", field.Name, err)
			}

			continue
		}

		values := entry.GetAttributeValues(ldapTag)
		if len(values) == 0 {
			continue
		}

		if err := setFieldValue(fieldValue, values); err != nil {
			return fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}

	return nil
}

// setFieldValue handles encoding.TextUnmarshaler and standard types with slice support.
func setFieldValue(fieldValue reflect.Value, values []string) error {
	fieldType := fieldValue.Type()

	// Handle slices
	if fieldType.Kind() == reflect.Slice {
		return setSliceValue(fieldValue, values)
	}

	// Single value - use first value
	if len(values) == 0 {
		return nil
	}

	return setSingleValue(fieldValue, values[0])
}

// setSliceValue handles slice types with TextUnmarshaler support.
func setSliceValue(fieldValue reflect.Value, values []string) error {
	elemType := fieldValue.Type().Elem()
	slice := reflect.MakeSlice(fieldValue.Type(), 0, len(values))

	for _, value := range values {
		elem := reflect.New(elemType).Elem()
		if err := setSingleValue(elem, value); err != nil {
			return fmt.Errorf("failed to set slice element: %w", err)
		}

		slice = reflect.Append(slice, elem)
	}

	fieldValue.Set(slice)

	return nil
}

// setSingleValue handles single values with TextUnmarshaler priority and standard type fallbacks.
func setSingleValue(fieldValue reflect.Value, value string) error {
	// Priority 1: Try encoding.TextUnmarshaler
	if fieldValue.CanAddr() {
		addr := fieldValue.Addr()
		if unmarshaler, ok := addr.Interface().(encoding.TextUnmarshaler); ok {
			return unmarshaler.UnmarshalText([]byte(value))
		}
	}

	// Priority 2: Handle standard types
	return setStandardType(fieldValue, value)
}

// setStandardType handles all Go standard types.
func setStandardType(fieldValue reflect.Value, value string) error {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)

		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldValue.Type() == reflect.TypeOf(time.Duration(0)) {
			// Handle time.Duration specially
			d, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("invalid duration: %w", err)
			}

			fieldValue.SetInt(int64(d))

			return nil
		}

		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer: %w", err)
		}

		fieldValue.SetInt(intVal)

		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid unsigned integer: %w", err)
		}

		fieldValue.SetUint(uintVal)

		return nil

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float: %w", err)
		}

		fieldValue.SetFloat(floatVal)

		return nil

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean: %w", err)
		}

		fieldValue.SetBool(boolVal)

		return nil

	case reflect.Struct:
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			// Handle time.Time
			timeVal, err := parseTimeValue(value)
			if err != nil {
				return fmt.Errorf("invalid time: %w", err)
			}

			fieldValue.Set(reflect.ValueOf(timeVal))

			return nil
		}

		return nil

	case reflect.Ptr:
		// Handle pointer types
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		}

		return setSingleValue(fieldValue.Elem(), value)

	default:
		return fmt.Errorf("unsupported field type: %v", fieldValue.Type())
	}
}

// parseTimeValue attempts multiple time formats.
func parseTimeValue(value string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"20060102150405Z", // LDAP timestamp format
	}

	for _, format := range formats {
		if t, err := time.Parse(format, value); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", value)
}

// Nullable types example.
type NullableInt struct {
	Value int
	Valid bool
}

func (n *NullableInt) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		n.Valid = false

		return nil
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return err
	}

	n.Value = val
	n.Valid = true

	return nil
}
