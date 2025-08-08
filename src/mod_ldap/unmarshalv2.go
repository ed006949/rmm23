package mod_ldap

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"time"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
)

// WalkTags processes struct fields with encoding.TextUnmarshaler support and standard type fallbacks.
func WalkTags(entries []*ldap.Entry, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("target must be pointer to slice")
	}

	sliceValue := targetValue.Elem()
	elemType := sliceValue.Type().Elem()

	for _, entry := range entries {
		item, err := createAndFillItem(entry, elemType)
		if err != nil {
			return fmt.Errorf("failed to process entry %s: %w", entry.DN, err)
		}

		sliceValue.Set(reflect.Append(sliceValue, item))
	}

	return nil
}

// createAndFillItem creates the appropriate item type and fills it
func createAndFillItem(entry *ldap.Entry, elemType reflect.Type) (reflect.Value, error) {
	var structValue reflect.Value
	var item reflect.Value

	if elemType.Kind() == reflect.Ptr {
		// Handle []*Entry
		structType := elemType.Elem()
		if structType.Kind() != reflect.Struct {
			return reflect.Value{}, fmt.Errorf("pointer element must point to struct, got %v", structType)
		}

		item = reflect.New(structType) // Create *Entry
		structValue = item.Elem()      // Get Entry for filling
	} else {
		// Handle []Entry
		if elemType.Kind() != reflect.Struct {
			return reflect.Value{}, fmt.Errorf("slice element must be struct, got %v", elemType)
		}

		item = reflect.New(elemType).Elem() // Create Entry
		structValue = item                  // Use directly for filling
	}

	if err := walkStructFields(entry, structValue); err != nil {
		return reflect.Value{}, err
	}

	return item, nil
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
		if ldapTag == "dn" { // just use entryDN attribute
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
	// Handle special types first before TextUnmarshaler check
	switch fieldValue.Type() {
	case reflect.TypeOf(time.Time{}):
		timeVal, err := parseTimeValue(value)
		if err != nil {
			return fmt.Errorf("invalid time: %w", err)
		}
		fieldValue.Set(reflect.ValueOf(timeVal))
		return nil

	case reflect.TypeOf(time.Duration(0)):
		d, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}
		fieldValue.SetInt(int64(d))
		return nil
	}

	// Priority: Try encoding.TextUnmarshaler
	if fieldValue.CanAddr() {
		addr := fieldValue.Addr()
		if unmarshaler, ok := addr.Interface().(encoding.TextUnmarshaler); ok {
			return unmarshaler.UnmarshalText([]byte(value))
		}
	}

	// Fallback: Handle standard types
	return setStandardType(fieldValue, value)
}

// setStandardType handles all Go standard types.
func setStandardType(fieldValue reflect.Value, value string) error {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
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

	case reflect.Ptr:
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		}
		return setSingleValue(fieldValue.Elem(), value)

	default:
		return fmt.Errorf("unsupported field type: %v", fieldValue.Type())
	}
}

// parseTimeValue with smart format detection based on string length and pattern
func parseTimeValue(value string) (t time.Time, err error) {
	// Skip empty value
	switch {
	case len(value) == 0:
		return time.Time{}, nil
	}

	// Try BER GeneralizedTime first (handles LDAP timestamps)
	switch t, err = ber.ParseGeneralizedTime([]byte(value)); {
	case err == nil:
		return
	}

	t = time.Time{}

	// Try time.Time's built-in UnmarshalText (handles RFC3339, etc.)
	switch err = t.UnmarshalText([]byte(value)); {
	case err == nil:
		return
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", value)
}
