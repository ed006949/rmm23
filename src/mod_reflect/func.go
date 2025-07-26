package mod_reflect

import (
	"reflect"

	"rmm23/src/mod_errors"
)

func GetStructRV(inbound any) (outboundRV reflect.Value, err error) {
	var (
		rv = reflect.ValueOf(inbound)
	)

	switch {
	case rv.Kind() == reflect.Ptr:
		rv = rv.Elem()
	}

	switch {
	case rv.Kind() != reflect.Struct:
		return outboundRV, mod_errors.ENotStructOrPtrStruct
	}

	return rv, nil
}

func GetStructRT(inbound any) (outboundRT reflect.Type, err error) {
	var (
		rt = reflect.TypeOf(inbound)
	)
	switch {
	case rt.Kind() == reflect.Ptr:
		rt = rt.Elem()
	}

	switch {
	case rt.Kind() != reflect.Struct:
		return outboundRT, mod_errors.ENotStructOrPtrStruct
	}

	return rt, nil
}

func GetStructSVST(inbound any) (outboundSV reflect.Value, outboundST reflect.Type, err error) {
	switch {
	case reflect.ValueOf(inbound).Kind() != reflect.Ptr:
		return outboundSV, outboundST, mod_errors.ENotPtr
	}

	var (
		sv, st = reflect.ValueOf(inbound).Elem(), reflect.TypeOf(inbound).Elem()
	)

	switch {
	case sv.Kind() != reflect.Struct:
		return outboundSV, outboundST, mod_errors.ENotStruct
	}

	return sv, st, nil
}
