package l

import (
	"net/url"
	"time"

	"github.com/rs/zerolog"
)

type Z map[string]interface{}

type nameType string
type configType string
type dryRunType string
type modeType string
type nodeType string
type dbType string
type verbosityType string
type gitCommitType string

type nameValue string
type configValue string
type dryRunValue bool
type modeValue int
type nodeValue int
type verbosityValue zerolog.Level
type commitValue string

type DaemonConfig struct {
	Name      nameValue      `env:"name,omitempty" json:"name,omitempty"`           //
	Verbosity verbosityValue `env:"verbosity,omitempty" json:"verbosity,omitempty"` //
	DryRun    dryRunValue    `env:"dry-run,omitempty" json:"dry-run,omitempty"`     //
	Mode      modeValue      `env:"mode,omitempty" json:"mode,omitempty"`           //
	Node      nodeValue      `env:"node,omitempty" json:"node,omitempty"`           //
	DB        *url.URL       `env:"db,omitempty" json:"db,omitempty"`               //
	Config    configValue    `env:"config,omitempty" json:"-"`                      //
	run       runType        //
	build     buildType      //
}

// run represents operational settings
type runType struct {
	name      string        //
	verbosity zerolog.Level //
	dryRun    bool          //
	mode      int           //
	node      int           //
	db        *url.URL      //
	config    string        //
	time      time.Time     //
	commit    string        //
}

// buildType represents buildType-time settings
type buildType struct {
	name      string //
	verbosity string //
	dryRun    string //
	mode      string //
	node      string //
	db        string //
	config    string //
	time      string //
	commit    string //
}
