package mod_strings

import (
	"fmt"
	"strings"
)

func (r EntryFieldName) String() (outbound string)         { return string(r) }
func (r EntryFieldName) FieldName() (outbound string)      { return JSONPathHeader + r.String() }
func (r EntryFieldName) FieldNameSlice() (outbound string) { return r.FieldName() + "[*]" }

func (r *EntryFieldMap) buildFVQuery(field EntryFieldName, value string) (outbound string) {
	return fmt.Sprintf(
		"@%s:%s%v%s",
		field.String(),
		FVEnclosure[(*r)[field]][0],
		escapeRedisQueryValue(value),
		FVEnclosure[(*r)[field]][1],
	)
}

func (r *EntryFieldMap) BuildQuery(inbound *FVs) (outbound string) {
	var (
		interim = make([]string, len(*r), len(*r))
	)

	for i, fv := range *inbound {
		interim[i] = r.buildFVQuery(fv.Field, fv.Value)
	}

	return strings.Join(interim, " ")
}
