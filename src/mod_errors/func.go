package mod_errors

import (
	"strings"
)

func StripErr(err error)                                 {}
func StripErr1[E any](inbound E, err error) (outbound E) { return inbound }

func PanicErr1[E any](inbound E, err error) (outbound E) {
	switch {
	case err != nil:
		panic(err)
	}
	return inbound
}

// Contains reports whether subErr is within err.
func Contains(err error, subErr error) bool {
	return err != nil && strings.Contains(err.Error(), subErr.Error())
}
