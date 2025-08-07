package mod_ldap

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-ldap/ldap/v3"
)

func UnmarshalEntry(e *ldap.Entry, out interface{}) error {
	vo := reflect.ValueOf(out)
	if vo.Kind() != reflect.Ptr || vo.IsNil() || vo.Elem().Kind() != reflect.Struct {
		return errors.New("UnmarshalEntry: expected pointer to struct")
	}

	val := vo.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fv := val.Field(i)

		if field.PkgPath != "" {
			continue // skip unexported
		}

		tag := field.Tag.Get("ldap")
		if tag == "" {
			continue
		}

		if tag == "dn" {
			if fv.Kind() == reflect.String && fv.CanSet() {
				fv.SetString(e.DN)
			}

			continue
		}

		values := getAttributeValues(e, tag)
		if len(values) == 0 {
			continue
		}

		// 1. Pointer to scalar
		if isPointerToScalar(fv) {
			if err := assignPointerToScalar(fv, values[0]); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}

			continue
		}
		// 2. Pointer to slice
		if isPointerToSlice(fv) {
			if err := assignPointerToSlice(fv, values); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}

			continue
		}
		// 3. Slice of pointers
		if isSliceOfPointers(fv) {
			if err := assignSliceOfPointers(fv, values); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}

			continue
		}
		// 4. Slice of values
		if fv.Kind() == reflect.Slice && fv.CanSet() {
			if err := assignSlice(fv, values); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}

			continue
		}
		// 5. Scalar
		if err := assignScalar(fv, values[0]); err != nil {
			return fmt.Errorf("%s: %w", field.Name, err)
		}
	}

	return nil
}

// ==== Field shape utilities ====

func isPointerToScalar(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && v.CanSet() && v.Type().Elem().Kind() != reflect.Slice
}

func isPointerToSlice(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && v.CanSet() && v.Type().Elem().Kind() == reflect.Slice
}

func isSliceOfPointers(v reflect.Value) bool {
	return v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Ptr && v.CanSet()
}

// ==== Helper assignment functions ====

func assignPointerToScalar(fv reflect.Value, val string) error {
	elemType := fv.Type().Elem()
	if fv.IsNil() {
		fv.Set(reflect.New(elemType))
	}

	elemPtr := fv
	if u, ok := elemPtr.Interface().(encoding.TextUnmarshaler); ok {
		return u.UnmarshalText([]byte(val))
	}

	elem := fv.Elem()

	return assignBasic(elem, val)
}

func assignPointerToSlice(fv reflect.Value, values []string) error {
	sliceType := fv.Type().Elem()

	elemType := sliceType.Elem()
	if fv.IsNil() {
		fv.Set(reflect.New(sliceType))
	}

	slice := reflect.MakeSlice(sliceType, len(values), len(values))
	for j, v := range values {
		elemPtr := reflect.New(elemType)
		if u, ok := elemPtr.Interface().(encoding.TextUnmarshaler); ok {
			if err := u.UnmarshalText([]byte(v)); err != nil {
				return err
			}

			slice.Index(j).Set(elemPtr.Elem())

			continue
		}

		if err := assignBasic(elemPtr.Elem(), v); err != nil {
			return err
		}
	}

	fv.Elem().Set(slice)

	return nil
}

func assignSliceOfPointers(fv reflect.Value, values []string) error {
	elemPtrType := fv.Type().Elem() // *T
	elemType := elemPtrType.Elem()  // T

	slice := reflect.MakeSlice(fv.Type(), len(values), len(values))
	for j, v := range values {
		elemPtr := reflect.New(elemType) // *T
		if u, ok := elemPtr.Interface().(encoding.TextUnmarshaler); ok {
			if err := u.UnmarshalText([]byte(v)); err != nil {
				return err
			}

			slice.Index(j).Set(elemPtr)

			continue
		}

		if err := assignBasic(elemPtr.Elem(), v); err != nil {
			return err
		}

		slice.Index(j).Set(elemPtr)
	}

	fv.Set(slice)

	return nil
}

func assignSlice(fv reflect.Value, values []string) error {
	elemType := fv.Type().Elem()

	slice := reflect.MakeSlice(fv.Type(), len(values), len(values))
	for j, v := range values {
		elemPtr := reflect.New(elemType)
		if u, ok := elemPtr.Interface().(encoding.TextUnmarshaler); ok {
			if err := u.UnmarshalText([]byte(v)); err != nil {
				return err
			}

			slice.Index(j).Set(elemPtr.Elem())

			continue
		}

		if err := assignBasic(slice.Index(j), v); err != nil {
			return err
		}
	}

	fv.Set(slice)

	return nil
}

func assignScalar(fv reflect.Value, val string) error {
	ptr := fv
	if fv.CanAddr() {
		ptr = fv.Addr()
	}

	if u, ok := ptr.Interface().(encoding.TextUnmarshaler); ok {
		return u.UnmarshalText([]byte(val))
	}

	return assignBasic(fv, val)
}

// ==== Centralized type conversion ====

func assignBasic(fv reflect.Value, val string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		fv.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}

		fv.SetUint(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}

		fv.SetBool(b)
	default:
		return fmt.Errorf("unsupported element type: %v", fv.Type().String())
	}

	return nil
}

// Helper: extract attribute values from github.com/go-ldap/ldap/v3 entries.
func getAttributeValues(e *ldap.Entry, name string) []string {
	for _, attr := range e.Attributes {
		if attr.Name == name {
			return attr.Values
		}
	}

	return nil
}
