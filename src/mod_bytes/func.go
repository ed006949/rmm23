package mod_bytes

import (
	"bytes"
	"fmt"
	"sort"

	"rmm23/src/mod_slices"
)

func JoinAny[S ~[]E, E any](inbound S, sep []byte, flag mod_slices.FlagType) (outbound []byte) {
	return bytes.Join(ToBytes(inbound, flag), sep)
}

func Join(inbound [][]byte, sep []byte, flag mod_slices.FlagType) (outbound []byte) {
	return bytes.Join(Normalize(inbound, flag), sep)
}

func FilterEmpty(inbound [][]byte) (outbound [][]byte) {
	var (
		a []byte
	)

	outbound = make([][]byte, 0, len(inbound))

	for _, b := range inbound {
		switch {
		case bytes.Compare(a, b) != 0:
			outbound = append(outbound, b)
		}
	}

	return
}

func TrimSpace(inbound [][]byte) (outbound [][]byte) {
	outbound = make([][]byte, 0, len(inbound))

	for a, b := range inbound {
		outbound[a] = bytes.TrimSpace(b)
	}

	return
}

func ToBytes[S ~[]E, E any](inbound S, flag mod_slices.FlagType) (outbound [][]byte) {
	outbound = make([][]byte, 0, len(inbound))

	for a, b := range inbound {
		outbound[a] = fmt.Append(nil, b)
	}

	outbound = Normalize(outbound, flag)

	return
}

func Split(inbound []byte, sep []byte, flag mod_slices.FlagType) (outbound [][]byte) {
	return Normalize(bytes.Split(inbound, sep), flag)
}

func SplitN(inbound []byte, sep []byte, n int, flag mod_slices.FlagType) (outbound [][]byte) {
	return Normalize(bytes.SplitN(inbound, sep, n), flag)
}

func Normalize(inbound [][]byte, flag mod_slices.FlagType) (outbound [][]byte) {
	switch {
	case flag.Has(mod_slices.FlagTrimSpace):
		inbound = TrimSpace(inbound)
	}

	switch {
	case flag.Has(mod_slices.FlagFilterEmpty):
		inbound = FilterEmpty(inbound)
	}

	switch {
	case flag.Has(mod_slices.FlagSort):
		Sort(inbound)
	}

	switch {
	case flag.Has(mod_slices.FlagCompact):
		inbound = Compact(inbound)
	}

	return inbound
}

func Sort(inbound [][]byte) {
	sort.Slice(inbound, func(i, j int) bool {
		return bytes.Compare(inbound[i], inbound[j]) < 0
	})
}

func Compact(inbound [][]byte) (outbound [][]byte) {
	switch len(inbound) {
	case 0, 1:
		return inbound
	}

	for k := 1; k < len(inbound); k++ {
		switch {
		case bytes.Equal(inbound[k], inbound[k-1]):
			var (
				s2 = inbound[k:]
			)
			for k2 := 1; k2 < len(s2); k2++ {
				switch {
				case !bytes.Equal(s2[k2], s2[k2-1]):
					inbound[k] = s2[k2]
					k++
				}
			}

			clear(inbound[k:])

			return inbound[:k]
		}
	}

	return inbound
}
