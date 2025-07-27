package mod_db

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
)
