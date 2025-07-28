package mod_db

import (
	"math"
	"time"
)

const (
	sliceSeparator byte = 0x1f
	jsonPathHeader      = "$."
	tagSeparator        = ","
)

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
	entryTypeEmpty attrEntryType = iota
	entryTypeDomain
	entryTypeGroup
	entryTypeUser
	entryTypeHost
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
	_ldap  = "ldap"
	_entry = "entry"
)

const (
	entryPrefix = "entry:"
)
const (
	entryDocIDHeader = _ldap + ":" + _entry + ":"
)

const (
	entryStatusUnknown attrEntryStatus = iota
	entryStatusLoaded
	entryStatusCreated
	entryStatusUpdated
	entryStatusDeleted
	entryStatusInvalid
	entryStatusReady = math.MaxInt
)
