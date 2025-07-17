package mod_db

const (
	redisTagName                = "redis"
	rediSearchTagName           = "redisearch"
	rediSearchTagTypeIgnore     = "-"
	rediSearchTagTypeTag        = "tag"
	rediSearchTagTypeText       = "text"
	rediSearchTagTypeNumeric    = "numeric"
	rediSearchTagOptionSortable = "sortable"
)

const (
	entryTypeDomain EntryType = iota
	entryTypeGroup
	entryTypeUser
	entryTypeHost
)
