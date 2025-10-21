package mod_db

import (
	"strings"
)

func escapeRedisQueryValue(inbound string) (outbound string) {
	replacer := strings.NewReplacer(
		`=`, `\=`, //
		`,`, `\,`, //
		`(`, `\(`, //
		`)`, `\)`, //
		`{`, `\{`, //
		`}`, `\}`, //
		`[`, `\[`, //
		`]`, `\]`, //
		`"`, `\"`, //
		`'`, `\'`, //
		`~`, `\~`, //
		`-`, `\-`, // (?)
	)

	return replacer.Replace(inbound)
}
