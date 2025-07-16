package mod_bools

import (
	"strings"

	"rmm23/src/mod_errors"
)

func StripIfBool1[E any](inbound E, flag bool) (outbound E) {
	switch {
	case flag:
		return inbound
	default:
		return
	}
}

func Parse(inbound string) (bool, error) {
	switch {
	case len(inbound) == 0:
		return false, mod_errors.ENODATA
	}
	inbound = strings.ToLower(inbound)

	switch inbound {
	case "1", "t", "y", "true", "yes", "on":
		return true, nil
	case "0", "f", "n", "false", "no", "off":
		return false, nil
	default:
		return false, mod_errors.EINVAL
	}
}
func FormatBool(inbound bool) string {
	switch inbound {
	case true:
		return "true"
	default:
		return "false"
	}
}
