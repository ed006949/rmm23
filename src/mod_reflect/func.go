package mod_reflect

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
)

func InitIfZero[T any](ptr *T) bool {
	var zeroValue T
	if reflect.DeepEqual(*ptr, zeroValue) {
		*ptr = zeroValue

		return true
	}

	return false
}

func NewPointerIfNil[T any](ptr **T) bool {
	if *ptr == nil {
		*ptr = new(T)

		return true
	}

	return false
}

func MakeMapIfNil[M ~map[K]V, K comparable, V any](m *M, size ...int) bool {
	switch {
	case *m != nil:
		return false
	}

	switch len(size) {
	case makeParam0:
		*m = make(M)

		return true
	case makeParam1:
		*m = make(M, size[0])

		return true
	default:
		panic(mod_errors.EUnwilling)
	}
}

func MakeSliceIfNil[S ~[]V, V any](s *S, size ...int) bool {
	switch {
	case *s != nil:
		return false
	}

	switch len(size) {
	case makeParam1:
		*s = make(S, size[0])

		return true
	case makeParam2:
		*s = make(S, size[0], size[1])

		return true
	default:
		panic(mod_errors.EUnwilling)
	}
}

func RetryWithCtx(ctx context.Context, maxTries int, interval time.Duration, fn func() error) (err error) {
	for attempt := 1; maxTries == 0 || attempt <= maxTries; attempt++ {
		switch err = fn(); {
		case err != nil:
			log.Warn().Int("attempt", attempt).Int("max", maxTries).Err(err).Msg("retry")

			switch err = WaitWithCtx(ctx, interval); {
			case err != nil:
				return
			}
		default:
			return
		}
	}

	return
}

func WaitWithCtx(ctx context.Context, d time.Duration) (err error) {
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
