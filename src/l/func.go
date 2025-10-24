package l

import (
	"strings"
)

// envName generates an environment variable name based on the daemon name and a daemon flag name.
// It converts the combined name to uppercase and replaces hyphens with underscores.
func envName(inbound int) (outbound string) {
	return strings.ReplaceAll(
		strings.ToUpper(
			buildName+"_"+daemonFlagName[inbound]),
		"-",
		"_",
	)
}
