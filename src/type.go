package main

import (
	"net/netip"

	"rmm23/src/l"
	"rmm23/src/mod_db"
	"rmm23/src/mod_ldap"
)

type ConfigRoot struct {
	Conf Conf `json:"conf,omitempty"`
}

type Conf struct {
	Daemon     *l.DaemonConfig `json:"daemon,omitempty"`
	DB         *mod_db.Conf    `json:"db,omitempty"`
	Networking *ConfNetworking `json:"networking,omitempty"`
	LDAP       *mod_ldap.Conf  `json:"ldap,omitempty"`
	Legacy     *ConfLegacy     `json:"legacy,omitempty"`
}

type ConfLegacy struct {
	PKI string `json:"PKI,omitempty"`
}
type ConfNetworking struct {
	User *ConfNetworkingUser `json:"user,omitempty"`
}
type ConfNetworkingUser struct {
	Subnet netip.Prefix `json:"subnet,omitempty"`
	Bits   int          `json:"bits,omitempty"`
}
