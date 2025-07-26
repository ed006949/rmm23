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
	entryTypeEmpty AttrEntryType = iota
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
	entryDocIDHeader = "ldap" + ":" + "entry" + ":"
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
	_type   entryFieldName = "type"
	_status entryFieldName = "status"
	_baseDN entryFieldName = "baseDN"

	_uuid            entryFieldName = "uuid"
	_dn              entryFieldName = "dn"
	_objectClass     entryFieldName = "objectClass"
	_creatorsName    entryFieldName = "creatorsName"
	_createTimestamp entryFieldName = "createTimestamp"
	_modifiersName   entryFieldName = "modifiersName"
	_modifyTimestamp entryFieldName = "modifyTimestamp"

	_cn                   entryFieldName = "cn"
	_dc                   entryFieldName = "dc"
	_description          entryFieldName = "description"
	_destinationIndicator entryFieldName = "destinationIndicator"
	_displayName          entryFieldName = "displayName"
	_gidNumber            entryFieldName = "gidNumber"
	_homeDirectory        entryFieldName = "homeDirectory"
	_ipHostNumber         entryFieldName = "ipHostNumber"
	_mail                 entryFieldName = "mail"
	_member               entryFieldName = "member"
	_o                    entryFieldName = "o"
	_ou                   entryFieldName = "ou"
	_owner                entryFieldName = "owner"
	_sn                   entryFieldName = "sn"
	_sshPublicKey         entryFieldName = "sshPublicKey"
	_telephoneNumber      entryFieldName = "telephoneNumber"
	_telexNumber          entryFieldName = "telexNumber"
	_uid                  entryFieldName = "uid"
	_uidNumber            entryFieldName = "uidNumber"
	_userPKCS12           entryFieldName = "userPKCS12"
	_userPassword         entryFieldName = "userPassword"

	_labeledURI entryFieldName = "labeledURI"
)
