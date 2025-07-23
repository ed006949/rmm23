package main

import (
	"encoding/xml"

	"rmm23/src/l"
	"rmm23/src/mod_ldap"
)

type ConfigRoot struct {
	Conf Conf `json:"conf"`
}

type Conf struct {
	Daemon l.DaemonConfig `json:"daemon"`
	LDAP   LDAPConfig     `json:"ldap"`
}

type LDAPConfig struct {
	URL      string        `json:"url"`
	Settings []LDAPSetting `json:"settings"`
	Domain   []LDAPDomain  `json:"domain"`
}

type LDAPSetting struct {
	Type   string `json:"type"`
	DN     string `json:"dn"`
	CN     string `json:"cn"`
	Filter string `json:"filter"`
}

type LDAPDomain struct {
	DN string `json:"dn"`
}

type xmlConf struct {
	XMLName xml.Name        `xml:"conf"`
	Daemon  *l.DaemonConfig `xml:"daemon,omitempty"`
	LDAP    *mod_ldap.Conf  `xml:"ldap,omitempty"`
}
