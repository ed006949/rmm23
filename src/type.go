package main

import (
	"rmm23/src/l"
	"rmm23/src/mod_ldap"
)

type ConfigRoot struct {
	Conf Conf `json:"conf"`
}

type Conf struct {
	Daemon l.DaemonConfig `json:"daemon"`
	LDAP   mod_ldap.LDAPConfig     `json:"ldap"`
}
