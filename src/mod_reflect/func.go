package mod_reflect

import (
	"context"
	"reflect"
	"strings"
	"time"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
)

func MakeMapIfNil[M ~map[K]V, K comparable, V any](m *M, size ...int) {
	switch {
	case *m == nil:
		switch len(size) {
		case makeParam0:
			*m = make(M)
		case makeParam1:
			*m = make(M, size[0])
		default:
		}
	}
}

func MakeSliceIfNil[S ~[]V, V any](s *S, size ...int) {
	switch {
	case *s == nil:
		switch len(size) {
		case makeParam1:
			*s = make(S, size[0])
		case makeParam2:
			*s = make(S, size[0], size[1])
		default:
		}
	}
}

func WaitCtx(ctx context.Context, d time.Duration) (err error) {
	var (
		t = time.NewTimer(d)
	)
	defer t.Stop()

	select {
	case <-t.C:
		return
	case <-ctx.Done():
		return ctx.Err()
	}
}

type FieldTypeInfo struct {
	Kind     reflect.Kind
	Type     reflect.Type // The field's own type
	ElemType reflect.Type // Non-nil if Kind == Slice
}

func BuildStructMap(inbound any, targetTag string) (outbound map[string]FieldTypeInfo, err error) {
	outbound = make(map[string]FieldTypeInfo)

	var (
		rt reflect.Type
	)
	switch rt, err = getStructRT(inbound); {
	case err != nil:
		return
	}

	for i := 0; i < rt.NumField(); i++ {
		var (
			field  = rt.Field(i)
			tagStr = field.Tag.Get(targetTag)
			tag    = strings.SplitN(tagStr, mod_strings.TagSeparator, mod_slices.KVElements)
		)
		switch {
		case len(tag) == 0 || len(tag[0]) == 0 || tag[0] == "-":
			continue
		}

		var (
			fTypeInfo = FieldTypeInfo{
				Kind: field.Type.Kind(),
				Type: field.Type,
			}
		)
		switch {
		case field.Type.Kind() == reflect.Slice:
			fTypeInfo.ElemType = field.Type.Elem()
		}

		outbound[tag[0]] = fTypeInfo
	}

	return
}

func getStructRT(inbound any) (outboundRT reflect.Type, err error) {
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

func getStructRV(inbound any) (outboundRV reflect.Value, err error) {
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

func getStructSVST(inbound any) (outboundSV reflect.Value, outboundST reflect.Type, err error) {
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
