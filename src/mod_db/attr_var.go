package mod_db

import (
	"strconv"
)

var (
	entryTypeName = map[attrEntryType]string{
		entryTypeEmpty:  "",
		entryTypeDomain: "domain",
		entryTypeGroup:  "group",
		entryTypeUser:   "user",
		entryTypeHost:   "host",
	}
	entryTypeID = map[string]attrEntryType{
		entryTypeName[entryTypeEmpty]:  entryTypeEmpty,
		entryTypeName[entryTypeDomain]: entryTypeDomain,
		entryTypeName[entryTypeGroup]:  entryTypeGroup,
		entryTypeName[entryTypeUser]:   entryTypeUser,
		entryTypeName[entryTypeHost]:   entryTypeHost,
	}
	entryTypeNumber = map[attrEntryType]string{
		entryTypeEmpty:  strconv.FormatInt(int64(entryTypeEmpty), 10),
		entryTypeDomain: strconv.FormatInt(int64(entryTypeDomain), 10),
		entryTypeGroup:  strconv.FormatInt(int64(entryTypeGroup), 10),
		entryTypeUser:   strconv.FormatInt(int64(entryTypeUser), 10),
		entryTypeHost:   strconv.FormatInt(int64(entryTypeHost), 10),
	}
)
