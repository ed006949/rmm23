package mod_slices

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

func FilterSlice[S ~[]E, E cmp.Ordered](inbound S, filter ...E) (outbound S) {
	var (
		interim = IndexSlice(filter)
	)
	for _, b := range inbound {
		switch _, ok := interim[b]; {
		case !ok:
			outbound = append(outbound, b)
		}
	}
	return
}
func IndexSlice[S ~[]E, E cmp.Ordered, M map[E]int](inbound S) (outbound M) {
	outbound = make(M)
	for a, b := range inbound {
		outbound[b] = a
	}
	return
}

func Sort[S ~[]E, E cmp.Ordered](inbound S)    { slices.Sort(inbound) }
func Compact[S ~[]E, E cmp.Ordered](inbound S) { slices.Compact(inbound) }

func Normalize[S ~[]E, E cmp.Ordered](inbound S) {
	Sort(inbound)
	Compact(inbound)
}

func Join[S ~[]E, E cmp.Ordered](inbound S, sep string) string {
	var (
		interim []string
	)
	for _, b := range inbound {
		switch d := fmt.Sprint(b); {
		case len(d) > 0:
			interim = append(interim, d)
		}
	}
	return strings.Join(interim, sep)
}

func NormalizeAndJoin[S ~[]E, E cmp.Ordered](inbound S, sep string) string {
	Normalize(inbound)
	return Join(inbound, sep)
}
