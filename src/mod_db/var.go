package mod_db

var (
	entryTypeName = map[AttrType]string{
		entryTypeEmpty:  "",
		entryTypeDomain: "domain",
		entryTypeGroup:  "group",
		entryTypeUser:   "user",
		entryTypeHost:   "host",
	}
	entryTypeID = map[string]AttrType{
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
