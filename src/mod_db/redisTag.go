package mod_db

const (
	redisTagName                = "redis"
	redisearchTagName           = "redisearch"
	redisearchTagTypeIgnore     = "-"
	redisearchTagTypeTag        = "TAG"     // exact match. string, bool.
	redisearchTagTypeGeo        = "GEO"     // ....
	redisearchTagTypeText       = "TEXT"    // partial match. string.
	redisearchTagTypeNumeric    = "NUMERIC" // numeric search. ints.
	redisearchTagOptionSortable = "sortable"
)
