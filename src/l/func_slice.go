package l

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strings"
)

func FilterSlice[S ~[]E, E comparable](inbound S, filter ...E) (outbound S) {
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
func IndexSlice[S ~[]E, E comparable, M map[E]int](inbound S) (outbound M) {
	outbound = make(M)
	for a, b := range inbound {
		outbound[b] = a
	}
	return
}
func NormalizeSlice[S ~[]E, E interface {
	comparable
	~string | ~[]byte
}, M map[E]struct{}](inbound S) (outbound S) {
	var (
		interimMap = make(M)
	)
	for _, b := range inbound {
		switch {
		case len(b) == 0:
			continue
		}
		switch _, ok := interimMap[b]; {
		case !ok:
			interimMap[b] = struct{}{}
			outbound = append(outbound, b)
		}
	}
	sort.Slice(outbound, func(i, j int) bool { return i < j })
	return
}
func JoinSlice[S ~[]E, E cmp.Ordered](inbound S, sep string) string {
	slices.Sort(inbound)
	slices.Compact(inbound)
	var (
		interim S
	)
	return strings.Join(inbound, sep)
}
func StringsJoin(elems []string, sep string) string {
	var (
		interim []string
	)
	for _, elem := range elems {
		switch {
		case len(elem) > 0:
			interim = append(interim, elem)
		}
	}
	return strings.Join(interim, sep)
}
func StringsJoinAndRemoveDuplicatesAndSort(elems []string, sep string) string {
	var (
		interim  = make(map[string]struct{})
		outbound []string
	)
	for _, elem := range elems {
		switch {
		case len(elem) == 0:
			continue
		}
		switch _, ok := interim[elem]; {
		case !ok:
			interim[elem] = struct{}{}
			outbound = append(outbound, elem)
		}
	}
	sort.Strings(outbound)
	return strings.Join(outbound, sep)
}
func StringsRemoveDuplicatesAndSort(elems []string) (outbound []string) {
	var (
		interim = make(map[string]struct{})
	)
	for _, elem := range elems {
		switch {
		case len(elem) == 0:
			continue
		}
		switch _, ok := interim[elem]; {
		case !ok:
			interim[elem] = struct{}{}
			outbound = append(outbound, elem)
		}
	}
	sort.Strings(outbound)
	return
}
