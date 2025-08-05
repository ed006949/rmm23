package mod_db

import (
	"rmm23/src/mod_errors"
)

//
// attrEntryType

func (r attrEntryType) String() (outbound string) { return entryTypeName[r] }
func (r attrEntryType) Number() (outbound string) { return entryTypeNumber[r] }
func (r *attrEntryType) Parse(inbound string) (err error) {
	switch value, ok := entryTypeID[inbound]; {
	case !ok:
		return mod_errors.EUnknownType
	default:
		*r = value

		return
	}
}
