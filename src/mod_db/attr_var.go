package mod_db

import (
	"strconv"
)

var (
	entryTypeName = map[attrEntryType]string{
		entryTypeEmpty:  "",
		EntryTypeDomain: "domain",
		EntryTypeGroup:  "group",
		EntryTypeUser:   "user",
		EntryTypeHost:   "host",
	}
	entryTypeID = map[string]attrEntryType{
		entryTypeName[entryTypeEmpty]:  entryTypeEmpty,
		entryTypeName[EntryTypeDomain]: EntryTypeDomain,
		entryTypeName[EntryTypeGroup]:  EntryTypeGroup,
		entryTypeName[EntryTypeUser]:   EntryTypeUser,
		entryTypeName[EntryTypeHost]:   EntryTypeHost,
	}
	entryTypeNumber = map[attrEntryType]string{
		entryTypeEmpty:  strconv.FormatInt(int64(entryTypeEmpty), 10),
		EntryTypeDomain: strconv.FormatInt(int64(EntryTypeDomain), 10),
		EntryTypeGroup:  strconv.FormatInt(int64(EntryTypeGroup), 10),
		EntryTypeUser:   strconv.FormatInt(int64(EntryTypeUser), 10),
		EntryTypeHost:   strconv.FormatInt(int64(EntryTypeHost), 10),
	}
)
