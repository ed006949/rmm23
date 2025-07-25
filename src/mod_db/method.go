package mod_db

import (
	"rmm23/src/mod_errors"
)

func (r *AttrType) Parse(inbound string) (err error) {
	switch value, ok := entryTypeID[inbound]; {
	case !ok:
		return mod_errors.EUnknownType
	default:
		*r = value

		return
	}
}
