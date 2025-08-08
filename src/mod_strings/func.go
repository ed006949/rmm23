package mod_strings

import (
	"fmt"
	"strings"

	"rmm23/src/mod_slices"
)

func JoinAny[S ~[]E, E any](inbound S, sep string, flag mod_slices.FlagType) (outbound string) {
	return strings.Join(ToStrings(inbound, flag), sep)
}

func Join(inbound []string, sep string, flag mod_slices.FlagType) (outbound string) {
	return strings.Join(Normalize(inbound, flag), sep)
}

func TrimSpace(inbound []string) (outbound []string) {
	outbound = make([]string, len(inbound), len(inbound))

	for a, b := range inbound {
		outbound[a] = strings.TrimSpace(b)
	}

	return
}

func ToStrings[S ~[]E, E any](inbound S, flag mod_slices.FlagType) (outbound []string) {
	outbound = make([]string, 0, len(inbound))

	for a, b := range inbound {
		outbound[a] = fmt.Sprint(b)
	}

	outbound = Normalize(outbound, flag)

	return
}

func Split(inbound string, sep string, flag mod_slices.FlagType) (outbound []string) {
	return Normalize(strings.Split(inbound, sep), flag)
}

func SplitN(inbound string, sep string, n int, flag mod_slices.FlagType) (outbound []string) {
	return Normalize(strings.SplitN(inbound, sep, n), flag)
}

func Normalize(inbound []string, flag mod_slices.FlagType) (outbound []string) {
	switch {
	case flag.Has(mod_slices.FlagTrimSpace):
		inbound = TrimSpace(inbound)
	}

	switch {
	case flag.Has(mod_slices.FlagFilterEmpty):
		inbound = mod_slices.FilterEmpty(inbound)
	}

	switch {
	case flag.Has(mod_slices.FlagSort):
		mod_slices.Sort(inbound)
	}

	switch {
	case flag.Has(mod_slices.FlagCompact):
		inbound = mod_slices.Compact(inbound)
	}

	return inbound
}
