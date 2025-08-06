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
func Normalize[S ~[]E, E cmp.Ordered](inbound S, flag flagType) (outbound S) {
	// switch {
	// case flag.has(FlagTrimSpace):
	// 	// inbound = TrimStrings(inbound)
	// }
	switch {
	case flag.has(FlagFilterEmpty):
		inbound = FilterEmpty(inbound)
	}

	switch {
	case flag.has(FlagSort):
		Sort(inbound)
	}

	switch {
	case flag.has(FlagCompact):
		inbound = Compact(inbound)
	}

	return inbound
}

// Join concatenates the elements of the inbound slice into a single string, separated by 'sep'.
// The slice is normalized before joining based on the provided flags.
func Join[S ~[]E, E any](inbound S, sep string, flag flagType) (outbound string) {
	return strings.Join(ToStrings(inbound, flag), sep)
}

func JoinStrings(inbound []string, sep string, flag flagType) (outbound string) {
	return strings.Join(StringsNormalize(inbound, flag), sep)
}

func TrimStrings(inbound []string) (outbound []string) {
	for _, b := range inbound {
		outbound = append(outbound, strings.TrimSpace(b))
	}

	return
}

// ToStrings converts the elements of the inbound slice to their string representations.
// The slice is normalized before conversion based on the provided flags.
func ToStrings[S ~[]E, E any](inbound S, flag flagType) (outbound []string) {
	for _, b := range inbound {
		outbound = append(outbound, fmt.Sprint(b))
	}

	outbound = StringsNormalize(outbound, flag)

	return
}

func SplitString(inbound string, sep string, flag flagType) (outbound []string) {
	return StringsNormalize(strings.Split(inbound, sep), flag)
}

func SplitStringN(inbound string, sep string, n int, flag flagType) (outbound []string) {
	return StringsNormalize(strings.SplitN(inbound, sep, n), flag)
}

func StringsNormalize(inbound []string, flag flagType) (outbound []string) {
	switch {
	case flag.has(FlagTrimSpace):
		inbound = TrimStrings(inbound)
	}

	switch {
	case flag.has(FlagFilterEmpty):
		inbound = FilterEmpty(inbound)
	}

	switch {
	case flag.has(FlagSort):
		Sort(inbound)
	}

	switch {
	case flag.has(FlagCompact):
		inbound = Compact(inbound)
	}

	return inbound
}
