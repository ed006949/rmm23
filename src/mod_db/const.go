package mod_db

const (
	redisTagName                = "redis"
	rediSearchTagName           = "redisearch"
	rediSearchTagTypeIgnore     = "-"
	rediSearchTagTypeText       = "text"
	rediSearchTagTypeNumeric    = "numeric"
	rediSearchTagOptionSortable = "sortable"
)

const (
	EntryTypeDomain EntryType = iota
	EntryTypeGroup
	EntryTypeUser
	EntryTypeHost
)
