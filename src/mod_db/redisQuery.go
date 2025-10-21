package mod_db

import (
	"fmt"
	"strings"

	"rmm23/src/mod_strings"
)

const (
	enclosureEmpty0  = ""
	enclosureEmpty1  = ""
	enclosureSquare0 = "["
	enclosureSquare1 = "]"
	enclosureCurly0  = "{"
	enclosureCurly1  = "}"
)

var (
	fvEnclosure = map[string][2]string{
		redisearchTagTypeText:    {enclosureEmpty0, enclosureEmpty1},
		redisearchTagTypeTag:     {enclosureCurly0, enclosureCurly1},
		redisearchTagTypeNumeric: {enclosureSquare0, enclosureSquare1},
		redisearchTagTypeGeo:     {enclosureSquare0, enclosureSquare1},
	}
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

func (r *ftInfoAttributes) buildFVQuery(field mod_strings.EntryFieldName, value string) (outbound string) {
	return fmt.Sprintf(
		"@%s:%s%v%s",
		field.String(),
		fvEnclosure[(*r)[field].Type][0],
		escapeRedisQueryValue(value),
		fvEnclosure[(*r)[field].Type][1],
	)
}

func (r *ftInfoAttributes) buildQuery(inbound *mod_strings.FVs) (outbound string) {
	var (
		interim = make([]string, len(*r), len(*r))
	)

	for i, fv := range *inbound {
		interim[i] = r.buildFVQuery(fv.Field, fv.Value)
	}

	return strings.Join(interim, " ")
}
