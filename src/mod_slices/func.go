package mod_slices

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// FilterSlice filters elements from a slice based on a variadic filter.
// It returns a new slice containing only the elements from 'inbound' that are NOT present in 'filter'.
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

// IndexSlice creates a map from slice elements to their indices.
// It takes a slice 'inbound' and returns a map where keys are the elements
// and values are their corresponding 0-based indices in the slice.
func IndexSlice[S ~[]E, E cmp.Ordered, M map[E]int](inbound S) (outbound M) {
	outbound = make(M)
	for a, b := range inbound {
		outbound[b] = a
	}
	return
}

// Sort sorts the elements of the slice 'inbound' in ascending order.
func Sort[S ~[]E, E cmp.Ordered](inbound S) { slices.Sort(inbound) }

// Compact removes consecutive duplicates from the slice 'inbound'.
// The order of the remaining elements is preserved.
func Compact[S ~[]E, E cmp.Ordered](inbound S) { slices.Compact(inbound) }

// Normalize sorts the slice 'inbound' and then removes consecutive duplicates.
// This effectively normalizes the slice by making its elements unique and ordered.
func Normalize[S ~[]E, E cmp.Ordered](inbound S) {
	Sort(inbound)
	Compact(inbound)
}

// Join concatenates the string representation of slice elements into a single string,
// separated by 'sep'. It converts each element to its string representation,
// trims whitespace, and only includes non-empty results.
func Join[S ~[]E, E any](inbound S, sep string) string {
	var (
		interim []string
	)
	for _, b := range inbound {
		// switch d := any(b).(type) {
		// case string:
		// 	switch f := strings.TrimSpace(d); {
		// 	case len(f) > 0:
		// 		interim = append(interim, f)
		// 	}
		// case fmt.Stringer:
		// 	switch f := strings.TrimSpace(d.String()); {
		// 	case len(f) > 0:
		// 		interim = append(interim, f)
		// 	}
		// default:
		// 	switch f := strings.TrimSpace(fmt.Sprint(d)); {
		// 	case len(f) > 0:
		// 		interim = append(interim, f)
		// 	}
		// }
		switch d := strings.TrimSpace(fmt.Sprint(b)); {
		case len(d) > 0:
			interim = append(interim, d)
		}
	}
	return strings.Join(interim, sep)
}

// NormalizeAndJoin normalizes the slice 'inbound' (sorts and compacts)
// and then joins its elements into a single string using 'sep'.
func NormalizeAndJoin[S ~[]E, E cmp.Ordered](inbound S, sep string) string {
	Normalize(inbound)
	return Join(inbound, sep)
}

// SplitAndNormalize splits the string representation of 'inbound' by 'sep',
// trims whitespace from the result, and then normalizes the resulting slice.
// Normalization includes sorting and removing duplicate elements.
func SplitAndNormalize[E any](inbound E, sep string) []string {
	var (
		interim = strings.Split(strings.TrimSpace(fmt.Sprint(inbound)), sep)
	)
	Normalize(interim)
	return interim
}
