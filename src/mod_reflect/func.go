package mod_reflect

import (
	"fmt"
	"reflect"
	"strings"

	"rmm23/src/mod_errors"
)

func WalkStructFields(structval interface{}, targetTag string) {
	t := reflect.TypeOf(structval)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	v := reflect.ValueOf(structval)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		kind := field.Type.Kind()
		jsonTag := field.Tag.Get(targetTag)

		// Format output with (json:"...") if tag is set
		tagInfo := ""
		if jsonTag != "" {
			tagInfo = fmt.Sprintf(" (%s:\"%s\")", targetTag, jsonTag)
		}

		switch kind {
		case reflect.Slice:
			fmt.Printf("%-24s : SLICE of %s%s\n", field.Name, field.Type.Elem(), tagInfo)
		case reflect.Map:
			fmt.Printf("%-24s : MAP (%s -> %s)%s\n", field.Name, field.Type.Key(), field.Type.Elem(), tagInfo)
		case reflect.Ptr:
			fmt.Printf("%-24s : POINTER to %s%s\n", field.Name, field.Type.Elem(), tagInfo)
		case reflect.Array:
			fmt.Printf("%-24s : ARRAY [%d] of %s%s\n", field.Name, field.Type.Len(), field.Type.Elem(), tagInfo)
		case reflect.Struct:
			fmt.Printf("%-24s : STRUCT (%s)%s\n", field.Name, field.Type.Name(), tagInfo)
		default:
			fmt.Printf("%-24s : VALUE (%s)%s\n", field.Name, kind, tagInfo)
		}
	}
}

func BuildStructMap(structval interface{}, targetTag string) (outboundKind map[string]reflect.Kind, outboundType map[string]reflect.Kind) {
	outboundKind = make(map[string]reflect.Kind)
	outboundType = make(map[string]reflect.Kind)

	var (
		v reflect.Value
		t reflect.Type
	)
	switch t = reflect.TypeOf(structval); {
	case t.Kind() == reflect.Ptr:
		t = t.Elem()
	}

	switch v = reflect.ValueOf(structval); {
	case v.Kind() == reflect.Ptr:
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		var (
			field = t.Field(i)
			kind  = field.Type.Kind()
			tag   = strings.Split(field.Tag.Get(targetTag), ",")
		)
		switch {
		case len(tag) == 0 || len(tag[0]) == 0:
			continue
		}

		outboundKind[tag[0]] = kind
		outboundType[tag[0]] = field.Type.Kind()
	}

	return
}

func CheckPointer(inbound reflect.Value) (outbound int) {
	var (
		kind     = inbound.Kind()
		elemKind = inbound.Type().Elem().Kind()
		canSet   = inbound.CanSet()
	)

	switch {
	case kind == reflect.Ptr && elemKind != reflect.Slice && canSet:
		return PointerToScalar
	case kind == reflect.Ptr && elemKind == reflect.Slice && canSet:
		return PointerToSlice
	case kind == reflect.Slice && elemKind == reflect.Ptr && canSet:
		return SliceOfPointers
	case kind == reflect.Slice && canSet:
		return SliceOfValues
	default:
		return
	}
}

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
