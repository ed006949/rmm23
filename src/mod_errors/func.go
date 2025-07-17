package mod_errors

import (
	"strings"
)

func StripErr(err error)                                 {}
func StripErr1[E any](inbound E, err error) (outbound E) { return inbound }

// Contains reports whether subErr is within err.
func Contains(err error, subErr error) bool { return strings.Contains(err.Error(), subErr.Error()) }
