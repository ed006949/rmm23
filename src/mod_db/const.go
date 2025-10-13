package mod_db

import (
	"time"
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
	// entryKeyHeader = _ldap + HeaderSeparator + _entry.
	certKeyHeader = _certificate
)
