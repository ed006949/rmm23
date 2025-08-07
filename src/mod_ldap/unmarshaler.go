package mod_ldap

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-ldap/ldap/v3"

	"rmm23/src/mod_reflect"
)

// Entry represents a single search result entry.
type Entry struct {
	// DN is the distinguished name of the entry
	DN string `json:"dn"`
	// Attributes are the returned attributes for the entry
	Attributes map[string][][]byte `json:"attributes"`
}

// // Entry represents a single search result entry
// type Entry struct {
// 	// DN is the distinguished name of the entry
// 	DN string `json:"dn"`
// 	// Attributes are the returned attributes for the entry
// 	Attributes []*ldap.EntryAttribute `json:"attributes"`
// }

// // EntryAttribute holds a single attribute
// type EntryAttribute struct {
// 	// Name is the name of the attribute
// 	Name string
// 	// Values contain the string values of the attribute
// 	Values []string
// 	// ByteValues contain the raw values of the attribute
// 	ByteValues [][]byte
// }

func UnmarshalEntry(e *ldap.Entry, out interface{}) error {
	var (
		vo = reflect.ValueOf(out)
	)
	switch {
	case vo.Kind() != reflect.Ptr || vo.IsNil() || vo.Elem().Kind() != reflect.Struct:
		return errors.New("UnmarshalEntry: expected pointer to struct")
	}

	var (
		val = vo.Elem()
		typ = val.Type()
	)

	for i := 0; i < typ.NumField(); i++ {
		var (
			field = typ.Field(i)
			fv    = val.Field(i)
		)
		switch {
		case len(field.PkgPath) != 0:
			continue // skip unexported fields
		}

		var (
			tag = field.Tag.Get(ldapTagName)
		)
		switch tag {
		case "":
			continue
		case "dn":
			switch {
			case fv.Kind() == reflect.String && fv.CanSet():
				fv.SetString(e.DN)
			}

			continue
		}

		var (
			values = getAttributeValues(e, tag)
		)
		switch {
		case len(values) == 0:
			continue
		}

		switch mod_reflect.CheckPointer(fv) {
		case mod_reflect.PointerToScalar: // 1. Pointer to scalar
			switch err := assignPointerToScalar(fv, values[0]); {
			case err != nil:
				return fmt.Errorf("%s: %w", field.Name, err)
			}
		case mod_reflect.PointerToSlice: // 2. Pointer to slice
			switch err := assignPointerToSlice(fv, values); {
			case err != nil:
				return fmt.Errorf("%s: %w", field.Name, err)
			}
		case mod_reflect.SliceOfPointers: // 3. Slice of pointers
			switch err := assignSliceOfPointers(fv, values); {
			case err != nil:
				return fmt.Errorf("%s: %w", field.Name, err)
			}
		case mod_reflect.SliceOfValues: // 4. Slice of values
			switch err := assignSlice(fv, values); {
			case err != nil:
				return fmt.Errorf("%s: %w", field.Name, err)
			}
		default: // 5. try Scalar
			switch err := assignScalar(fv, values[0]); {
			case err != nil:
				return fmt.Errorf("%s: %w", field.Name, err)
			}
		}
	}

	return nil
}

// ==== Helper assignment functions ====

func assignPointerToScalar(fv reflect.Value, val string) error {
	var (
		elemType = fv.Type().Elem()
	)
	switch {
	case fv.IsNil():
		fv.Set(reflect.New(elemType))
	}

	var (
		elemPtr = fv
	)
	switch u := elemPtr.Interface().(type) {
	case encoding.TextUnmarshaler:
		return u.UnmarshalText([]byte(val))
	}

	var (
		elem = fv.Elem()
	)

	return assignBasic(elem, val)
}

func assignPointerToSlice(fv reflect.Value, values []string) error {
	var (
		sliceType = fv.Type().Elem()
		elemType  = sliceType.Elem()
	)
	switch {
	case fv.IsNil():
		fv.Set(reflect.New(sliceType))
	}

	var (
		slice = reflect.MakeSlice(sliceType, len(values), len(values))
	)

	for j := 0; j < len(values); j++ {
		var (
			v       = values[j]
			elemPtr = reflect.New(elemType)
		)
		switch u := elemPtr.Interface().(type) {
		case encoding.TextUnmarshaler:
			switch err := u.UnmarshalText([]byte(v)); {
			case err != nil:
				return err
			}

			slice.Index(j).Set(elemPtr.Elem())

			continue
		}

		switch err := assignBasic(elemPtr.Elem(), v); {
		case err != nil:
			return err
		}

		slice.Index(j).Set(elemPtr.Elem())
	}

	fv.Elem().Set(slice)

	return nil
}

func assignSliceOfPointers(fv reflect.Value, values []string) error {
	var (
		elemPtrType = fv.Type().Elem()
		elemType    = elemPtrType.Elem()
		slice       = reflect.MakeSlice(fv.Type(), len(values), len(values))
	)
	for j := 0; j < len(values); j++ {
		var (
			v       = values[j]
			elemPtr = reflect.New(elemType) // *T
		)
		switch u := elemPtr.Interface().(type) {
		case encoding.TextUnmarshaler:
			switch err := u.UnmarshalText([]byte(v)); {
			case err != nil:
				return err
			}

			slice.Index(j).Set(elemPtr)

			continue
		}

		switch err := assignBasic(elemPtr.Elem(), v); {
		case err != nil:
			return err
		}

		slice.Index(j).Set(elemPtr)
	}

	fv.Set(slice)

	return nil
}

func assignSlice(fv reflect.Value, values []string) error {
	var (
		elemType = fv.Type().Elem()
		slice    = reflect.MakeSlice(fv.Type(), len(values), len(values))
	)
	for j := 0; j < len(values); j++ {
		var (
			v       = values[j]
			elemPtr = reflect.New(elemType)
		)
		switch u := elemPtr.Interface().(type) {
		case encoding.TextUnmarshaler:
			switch err := u.UnmarshalText([]byte(v)); {
			case err != nil:
				return err
			}

			slice.Index(j).Set(elemPtr.Elem())

			continue
		}

		switch err := assignBasic(slice.Index(j), v); {
		case err != nil:
			return err
		}
	}

	fv.Set(slice)

	return nil
}

func assignScalar(fv reflect.Value, val string) error {
	var ptr = fv
	switch {
	case fv.CanAddr():
		ptr = fv.Addr()
	}

	switch u := ptr.Interface().(type) {
	case encoding.TextUnmarshaler:
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
		var (
			n   int64
			err error
		)
		switch n, err = strconv.ParseInt(val, 10, 64); {
		case err != nil:
			return err
		}

		fv.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var (
			n   uint64
			err error
		)
		switch n, err = strconv.ParseUint(val, 10, 64); {
		case err != nil:
			return err
		}

		fv.SetUint(n)
	case reflect.Bool:
		var (
			b   bool
			err error
		)
		switch b, err = strconv.ParseBool(val); {
		case err != nil:
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
		switch {
		case attr.Name == name:
			return attr.Values
		}
	}

	return nil
}
