package mod_slices

import (
	"cmp"
	"slices"
)

func Filter[S ~[]E, E comparable](inbound S, filter ...E) (outbound S) {
	var (
		filters = Index(filter)
	)

	outbound = make(S, 0, len(inbound))

	for _, b := range inbound {
		switch _, ok := filters[b]; {
		case !ok:
			outbound = append(outbound, b)
		}
	}

	return
}

func Index[S ~[]E, E comparable, M map[E]struct{}](inbound S) (outbound M) {
	outbound = make(M)
	for _, b := range inbound {
		outbound[b] = struct{}{}
	}

	return
}
func HasIndex[T any](s []T, n int) bool {
	return n >= 0 && n < len(s)
}

func Sort[S ~[]E, E cmp.Ordered](inbound S) { slices.Sort(inbound) }

func Compact[S ~[]E, E comparable](inbound S) (outbound S) { return slices.Compact(inbound) }

func FilterEmpty[S ~[]E, E comparable](inbound S) (outbound S) {
	var (
		a E
	)

	return Filter(inbound, a)
}

func Normalize[S ~[]E, E cmp.Ordered](inbound S, flag FlagType) (outbound S) {
	// switch {
	// case flag.Has(FlagTrimSpace):
	// 	// inbound = TrimStrings(inbound)
	// }
	switch {
	case flag.Has(FlagFilterEmpty):
		inbound = FilterEmpty(inbound)
	}

	switch {
	case flag.Has(FlagSort):
		Sort(inbound)
	}

	switch {
	case flag.Has(FlagCompact):
		inbound = Compact(inbound)
	}

	return inbound
}
