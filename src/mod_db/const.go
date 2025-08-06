package mod_db

import (
	"math"
	"time"
)

const (
	sliceSeparator  = "\x1f"
	jsonPathHeader  = "$."
	tagSeparator    = ","
	headerSeparator = ":"
)

const (
	redisTagName                = "redis"
	redisearchTagName           = "redisearch"
	redisearchTagTypeIgnore     = "-"
	redisearchTagTypeTag        = "tag"
	redisearchTagTypeGeo        = "geo"
	redisearchTagTypeText       = "text"
	redisearchTagTypeNumeric    = "numeric"
	redisearchTagOptionSortable = "sortable"
)

const (
	entryTypeEmpty attrEntryType = iota
	EntryTypeDomain
	EntryTypeGroup
	EntryTypeUser
	EntryTypeHost
)

const (
	_PING = "PING"
)
const (
	connMaxIdle         = 4
	connMaxActive       = 4
	connIdleTimeout     = 240 * time.Second
	connWait            = true
	connMaxConnLifetime = 0
	connMaxPaging       = 1000000
)

const (
	_ldap        = "ldap"
	_entry       = "entry"
	_certificate = "certificate"
)

const (
	entryKeyHeader = _entry
	// entryKeyHeader = _ldap + headerSeparator + _entry.
	certKeyHeader = _certificate
)

const (
	entryStatusUnknown attrEntryStatus = iota
	entryStatusLoaded
	entryStatusCreated
	entryStatusUpdated
	entryStatusDeleted
	entryStatusInvalid
	entryStatusParsed
	entryStatusSanitized
	entryStatusReady = math.MaxInt
)
