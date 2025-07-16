package mod_slices

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

func isFlag(flag flag, flags ...flag) bool {
	for _, b := range flags {
		switch b {
		case FlagNormalize, flag:
			return true
		default:
		}
	}
	return false
}

// Filter filters elements from a slice based on a variadic filter.
// It returns a new slice containing only the elements from 'inbound' that are NOT present in 'filter'.
func Filter[S ~[]E, E cmp.Ordered](inbound S, filter ...E) {
	var (
		filters = Index(filter)
		interim S
	)
	for _, b := range inbound {
		switch _, ok := filters[b]; {
		case !ok:
			interim = append(interim, b)
		}
	}
	inbound = interim
}

func FilterEmpty[S ~[]E, E cmp.Ordered](inbound S) {
	var (
		// interim S
		a E
	)
	Filter(inbound, a)
	// for _, b := range inbound {
	// 	switch d := strings.TrimSpace(fmt.Sprint(b)); {
	// 	case len(d) > 0:
	// 		interim = append(interim, b)
	// 	}
	// }
	// inbound = interim
}

// Index creates a map from slice elements to their indices.
// It takes a slice 'inbound' and returns a map where keys are the elements
// and values are their corresponding 0-based indices in the slice.
func Index[S ~[]E, E cmp.Ordered, M map[E]int](inbound S) (outbound M) {
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
func Normalize[S ~[]E, E cmp.Ordered](inbound S, flags ...flag) {
	switch {
	case isFlag(FlagSort, flags...):
		Sort(inbound)
		fallthrough
	case isFlag(FlagCompact, flags...):
		Compact(inbound)
		fallthrough
	case isFlag(FlagFilterEmpty, flags...):
		FilterEmpty(inbound)
		// fallthrough
	}
}

// Join concatenates the string representation of slice elements into a single string,
// separated by 'sep'. It converts each element to its string representation,
// trims whitespace, and only includes non-empty results.
func Join[S ~[]E, E cmp.Ordered](inbound S, sep string, flags ...flag) string {
	return strings.Join(Strings(inbound, flags...), sep)
}

func Strings[S ~[]E, E cmp.Ordered](inbound S, flags ...flag) (outbound []string) {
	Normalize(inbound, flags...)
	for _, b := range inbound {
		outbound = append(outbound, fmt.Sprint(b))
	}
	return
}

// Split splits the string representation of 'inbound' by 'sep' and trims whitespace from the result.
func Split[E cmp.Ordered](inbound E, sep string) []string {
	return strings.Split(strings.TrimSpace(fmt.Sprint(inbound)), sep)
}

// SplitAndNormalize splits the string representation of 'inbound' by 'sep',
// trims whitespace from the result, and then normalizes the resulting slice.
// Normalization includes sorting and removing duplicate elements.
func SplitAndNormalize[E cmp.Ordered](inbound E, sep string) (outbound []string) {
	outbound = Split(inbound, sep)
	Normalize(outbound)
	return
}
