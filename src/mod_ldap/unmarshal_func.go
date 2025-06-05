package io_ldap

import (
	"reflect"
	"strings"

	"github.com/go-ldap/ldap/v3"

	"rmm23/src/l"
)

func readTag(f reflect.StructField) (options string, flag bool) {
	var (
		val, ok = f.Tag.Lookup(decoderTagName)
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
		flag = opts[1] == "omitempty"
	}
	return opts[0], flag
} // portions taken from "github.com/go-ldap/ldap/v3"

func unmarshal(e *ldap.Entry, i interface{}) (err error) {
	var (
		vo = reflect.ValueOf(i).Kind()
	)
	switch {
	case vo != reflect.Ptr:
		return ENotPtr
	}

	var (
		sv, st = reflect.ValueOf(i).Elem(), reflect.TypeOf(i).Elem()
	)
	switch {
	case sv.Kind() != reflect.Struct:
		return ENotStruct
	}

	for n := 0; n < st.NumField(); n++ {
		// Holds struct field value and type
		var (
			fv, ft = sv.Field(n), st.Field(n)
		)

		// skip unexported fields
		switch {
		case ft.PkgPath != "":
			continue
		}

		// omitempty can be safely discarded, as it's not needed when unmarshalling
		fieldTag, _ := readTag(ft)

		// Fill the field with the distinguishedName if the tag key is `dn`
		switch {
		case fieldTag == "dn":
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
} // portions taken from "github.com/go-ldap/ldap/v3"
