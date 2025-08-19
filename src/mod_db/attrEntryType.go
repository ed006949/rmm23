package mod_db

import (
	"strconv"

	"rmm23/src/mod_errors"
)

var (
	entryTypeName = map[attrEntryType]string{
		entryTypeEmpty:  "",
		entryTypeDomain: "domain",
		entryTypeGroup:  "group",
		entryTypeUser:   "user",
		EntryTypeHost:   "host",
	}
	entryTypeID = map[string]attrEntryType{
		entryTypeName[entryTypeEmpty]:  entryTypeEmpty,
		entryTypeName[entryTypeDomain]: entryTypeDomain,
		entryTypeName[entryTypeGroup]:  entryTypeGroup,
		entryTypeName[entryTypeUser]:   entryTypeUser,
		entryTypeName[EntryTypeHost]:   EntryTypeHost,
	}
	entryTypeNumber = map[attrEntryType]string{
		entryTypeEmpty:  strconv.FormatInt(int64(entryTypeEmpty), 10),
		entryTypeDomain: strconv.FormatInt(int64(entryTypeDomain), 10),
		entryTypeGroup:  strconv.FormatInt(int64(entryTypeGroup), 10),
		entryTypeUser:   strconv.FormatInt(int64(entryTypeUser), 10),
		EntryTypeHost:   strconv.FormatInt(int64(EntryTypeHost), 10),
	}
)

type attrEntryType int   //
type attrEntryStatus int //

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
