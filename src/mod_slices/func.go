// Package mod_slices provides utility functions for working with slices.
package mod_slices

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// Filter returns a new slice containing elements from 'inbound' that are not present in 'filter'.
func Filter[S ~[]E, E cmp.Ordered](inbound S, filter ...E) (outbound S) {
	var (
		filters = Index(filter)
	)

	for _, b := range inbound {
		switch _, ok := filters[b]; {
		case !ok:
			outbound = append(outbound, b)
		}
	}

	return
}

// Index creates a map where keys are elements from the inbound slice and values are their indices.
func Index[S ~[]E, E cmp.Ordered, M map[E]int](inbound S) (outbound M) {
	outbound = make(M)
	for a, b := range inbound {
		outbound[b] = a
	}

	return
}

// Sort sorts the elements of the inbound slice in ascending order.
func Sort[S ~[]E, E cmp.Ordered](inbound S) { slices.Sort(inbound) }

// Compact removes consecutive duplicate elements from the inbound slice.
func Compact[S ~[]E, E cmp.Ordered](inbound S) (outbound S) { return slices.Compact(inbound) }

// FilterEmpty filters out the zero value of type E from the inbound slice.
func FilterEmpty[S ~[]E, E cmp.Ordered](inbound S) (outbound S) {
	var (
		a E
	)

	return Filter(inbound, a)
}

// Normalize applies a series of transformations (sort, compact, filter empty) to the inbound slice based on the provided flags.
func Normalize[S ~[]E, E cmp.Ordered](inbound S, flag flag) (outbound S) {
	switch {
	case flag.has(FlagSort):
		Sort(inbound)

		fallthrough
	case flag.has(FlagCompact):
		inbound = Compact(inbound)

		fallthrough
	case flag.has(FlagFilterEmpty):
		inbound = FilterEmpty(inbound)
	}

	return inbound
}

// Join concatenates the elements of the inbound slice into a single string, separated by 'sep'.
// The slice is normalized before joining based on the provided flags.
func Join[S ~[]E, E cmp.Ordered](inbound S, sep string, flag flag) (outbound string) {
	return strings.Join(ToStrings(inbound, flag), sep)
}

// ToStrings converts the elements of the inbound slice to their string representations.
// The slice is normalized before conversion based on the provided flags.
func ToStrings[S ~[]E, E cmp.Ordered](inbound S, flag flag) (outbound []string) {
	var (
		converter func(in E) (out string)
	)

	switch {
	case flag.has(FlagTrimSpace):
		converter = func(in E) (out string) {
			return strings.TrimSpace(fmt.Sprint(in))
		}
	default:
		converter = func(in E) (out string) {
			return fmt.Sprint(in)
		}
	}

	// inbound = Normalize(inbound, flag)

	for _, b := range inbound {
		outbound = append(outbound, converter(b))
	}

	outbound = Normalize(outbound, flag)

	return
}

func Split(inbound string, sep string, flag flag) (outbound []string) {
	return Normalize(strings.Split(inbound, sep), flag)
}
