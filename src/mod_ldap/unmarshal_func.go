package mod_ldap

import (
	"reflect"
	"strings"

	"github.com/go-ldap/ldap/v3"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
)

// portions taken from "github.com/go-ldap/ldap/v3"

func readTag(f reflect.StructField) (options string, flag bool) {
	var (
		val, ok = f.Tag.Lookup(ldapTagName)
	)

	switch {
	case !ok:
		return f.Name, false
	}

	var (
		opts = strings.Split(val, ",")
	)

	switch {
	case len(opts) == 2:
		flag = opts[1] == ldapTagOptionOmitEmpty
	}

	return opts[0], flag
}

func UnmarshalEntry(e *ldap.Entry, i interface{}) (err error) {
	var (
		vo = reflect.ValueOf(i).Kind()
	)

	switch {
	case vo != reflect.Ptr:
		return mod_errors.ENotPtr
	}

	var (
		sv, st = reflect.ValueOf(i).Elem(), reflect.TypeOf(i).Elem()
	)

	switch {
	case sv.Kind() != reflect.Struct:
		return mod_errors.ENotStruct
	}

	for n := 0; n < st.NumField(); n++ {
		var (
			fv, ft = sv.Field(n), st.Field(n) // Holds struct field value and type
		)

		switch {
		case ft.PkgPath != "": // skip unexported fields
			continue
		}

		// omitempty can be safely discarded, as it's not needed when unmarshalling
		fieldTag, _ := readTag(ft)

		// Fill the field with the distinguishedName if the tag key is `dn`
		switch fieldTag {
		case "dn":
			switch _, err = ldap.ParseDN(e.DN); {
			case err != nil:
				return
			}

			fv.SetString(e.DN)

			continue
		}

		var (
			values = e.GetAttributeValues(fieldTag)
		)

		switch {
		case len(values) == 0:
			continue
		}

		switch fieldType := fv.Interface().(type) {
		default:
			switch reflect.TypeOf(fieldType).Kind() {
			case reflect.Map:
				var (
					ptrVal = reflect.MakeMap(reflect.TypeOf(fieldType))
				)

				switch unmarshaler, ok := ptrVal.Interface().(LDAPAttributeUnmarshaler); {
				case ok:
					switch err = unmarshaler.UnmarshalLDAPAttr(values); {
					case err != nil:
						l.Z{l.E: err, l.M: "LDAP Unmarshal", "DN": e.DN}.Warning()
						err = nil

						continue
					}

					fv.Set(ptrVal)
				}
			default:
				var (
					ptrVal = reflect.New(reflect.TypeOf(fieldType))
				)

				switch unmarshaler, ok := ptrVal.Interface().(LDAPAttributeUnmarshaler); {
				case ok:
					switch err = unmarshaler.UnmarshalLDAPAttr(values); {
					case err != nil:
						l.Z{l.E: err, l.M: "LDAP Unmarshal", "DN": e.DN}.Warning()
						err = nil

						continue
					}

					fv.Set(ptrVal.Elem())
				}
			}
		}
	}

	return
}
