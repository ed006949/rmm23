package main

import (
	"encoding/xml"

	"rmm23/src/l"
	"rmm23/src/mod_ldap"
)

type xmlConf struct {
	XMLName xml.Name       `xml:"conf"`
	Daemon  *l.ControlType `xml:"daemon,omitempty"`
	LDAP    *mod_ldap.Conf `xml:"ldap,omitempty"`
}
