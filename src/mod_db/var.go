package mod_db

var (
	entryTypeName = map[AttrEntryType]string{
		entryTypeEmpty:  "",
		entryTypeDomain: "domain",
		entryTypeGroup:  "group",
		entryTypeUser:   "user",
		entryTypeHost:   "host",
	}
	entryTypeID = map[string]AttrEntryType{
		entryTypeName[entryTypeEmpty]:  entryTypeEmpty,
		entryTypeName[entryTypeDomain]: entryTypeDomain,
		entryTypeName[entryTypeGroup]:  entryTypeGroup,
		entryTypeName[entryTypeUser]:   entryTypeUser,
		entryTypeName[entryTypeHost]:   entryTypeHost,
	}
)

// var (
// 	docFieldName = map[string]string{}
// )
