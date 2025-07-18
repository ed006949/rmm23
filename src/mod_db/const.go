package mod_db

const (
	redisTagName                = "redis"
	rediSearchTagName           = "redisearch"
	rediSearchTagTypeIgnore     = "-"
	rediSearchTagTypeTag        = "tag"
	rediSearchTagTypeGeo        = "geo"
	rediSearchTagTypeText       = "text"
	rediSearchTagTypeNumeric    = "numeric"
	rediSearchTagOptionSortable = "sortable"
)

const (
	entryTypeEmpty EntryType = iota
	entryTypeDomain
	entryTypeGroup
	entryTypeUser
	entryTypeHost
)
