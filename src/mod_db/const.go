package mod_db

import (
	"math"
	"time"
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
	entryTypeEmpty AttrType = iota
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
	entryDocIDHeader = "ldap" + ":" + "entry"
)

const (
	entryStatusUnknown AttrEntryStatus = iota
	entryStatusLoaded
	entryStatusCreated
	entryStatusUpdated
	entryStatusDeleted
	entryStatusInvalid
	entryStatusReady = math.MaxInt
)

const (
	_type   = "type"
	_baseDN = "baseDN"
	_uuid   = "uuid"
	_dn     = "dn"
)
