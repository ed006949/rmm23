package mod_ldap

import (
	"reflect"
	"strings"

	"github.com/go-ldap/ldap/v3"

	"rmm23/src/mod_reflect"
	"rmm23/src/mod_slices"
)

func readTag(f reflect.StructField) (options string, flag bool) {
	var (
		val, ok = f.Tag.Lookup(ldapTagName)
	)

	switch {
	case !ok:
		return f.Name, false
	}

	var (
		// opts = mod_slices.SplitString(val, mod_strings.TagSeparator, mod_slices.FlagNormalize)
		opts = strings.Split(val, ",")
	)

	switch {
	case len(opts) == mod_slices.KVElements:
		flag = opts[1] == ldapTagOptionOmitEmpty
	}

	return opts[0], flag
}

func UnmarshalEntry(e *ldap.Entry, i interface{}) (err error) {
	var (
		sv reflect.Value
		st reflect.Type
	)

	switch sv, st, err = mod_reflect.GetStructSVST(i); {
	case err != nil:
		return
	}

	for n := 0; n < st.NumField(); n++ {
		var (
			fv = sv.Field(n) // Holds struct field value
			ft = st.Field(n) // Holds struct field type
		)

		switch {
		case len(ft.PkgPath) != 0: // skip unexported fields
			continue
		}

		// omitempty can be safely discarded, as it's not needed when unmarshalling
		var (
			fieldTag, _ = readTag(ft)
		)

		// Fill the field with the distinguishedName if the tag key is `dn`
		switch fieldTag {
		case "dn":
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
			switch rt := reflect.TypeOf(fieldType); rt.Kind() {
			case reflect.Map:
				var (
					ptrVal = reflect.MakeMap(rt)
				)

				switch unmarshaler, ok := ptrVal.Interface().(LDAPAttributeUnmarshaler); {
				case ok:
					switch err = unmarshaler.UnmarshalLDAPAttr(values); {
					case err != nil:
						return
					}

					fv.Set(ptrVal)
				}
			default:
				var (
					ptrVal = reflect.New(rt)
				)

				switch unmarshaler, ok := ptrVal.Interface().(LDAPAttributeUnmarshaler); {
				case ok:
					switch err = unmarshaler.UnmarshalLDAPAttr(values); {
					case err != nil:
						return
					}

					fv.Set(ptrVal.Elem())
				}
			}
		}
	}

	return
}
