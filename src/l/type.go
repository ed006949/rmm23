package l

import (
	"net/url"

	"github.com/rs/zerolog"
)

type Z map[string]interface{}

type nameType string
type configType string
type dryRunType string
type modeType string
type verbosityType string
type gitCommitType string

type nameValue string
type configValue string
type dryRunFlag bool
type modeValue int
type verbosityLevel zerolog.Level
type gitCommitValue string

type ConfigRoot struct {
	Conf Conf `json:"conf"`
}

type Conf struct {
	Daemon DaemonConfig `json:"daemon"`
	Ldap   LDAPConfig   `json:"ldap"`
}

type DaemonConfig struct {
	Name      string         `json:"name"`
	Verbosity string         `json:"verbosity"`
	DryRun    bool           `json:"dry-run"`
	Node      int            `json:"node"`
	DB        url.URL        `json:"db"`
	GitCommit gitCommitValue `json:"-"`
	Config    configValue    `json:"-"`
}

type LDAPConfig struct {
	URL      string        `json:"url"`
	Settings []LDAPSetting `json:"settings"`
	Domain   []LDAPDomain  `json:"domain"`
}

type LDAPSetting struct {
	Type   string `json:"type"`
	Dn     string `json:"dn"`
	Cn     string `json:"cn"`
	Filter string `json:"filter"`
}

type LDAPDomain struct {
	Dn string `json:"dn"`
}

type ControlType struct {
	Name      nameValue      `xml:"name,attr,omitempty" json:"name,omitempty"`
	Config    configValue    `xml:"config,attr,omitempty" json:"config,omitempty"`
	DryRun    dryRunFlag     `xml:"dry-run,attr,omitempty" json:"dry-run,omitempty"`
	Mode      modeValue      `xml:"mode,attr,omitempty" json:"mode,omitempty"`
	Verbosity verbosityLevel `xml:"verbosity,attr,omitempty" json:"verbosity,omitempty"`
	GitCommit gitCommitValue `xml:"-" json:"-"`
}
