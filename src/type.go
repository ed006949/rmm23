package main

import (
	"rmm23/src/l"
	"rmm23/src/mod_db"
	"rmm23/src/mod_ldap"
)

type ConfigRoot struct {
	Conf Conf `json:"conf,omitempty"`
}

type Conf struct {
	Daemon *l.DaemonConfig `json:"daemon,omitempty"`
	DB     *mod_db.Conf    `json:"db,omitempty"`
	LDAP   *mod_ldap.Conf  `json:"ldap,omitempty"`
}
