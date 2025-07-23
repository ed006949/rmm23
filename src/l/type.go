package l

import (
	"net/url"
	"time"

	"github.com/rs/zerolog"
)

type Z map[string]interface{}

type DaemonConfig struct {
	Verbosity zerolog.Level `json:"verbosity,omitempty"` //
	DryRun    bool          `json:"dry-run,omitempty"`   //
	Mode      int           `json:"mode,omitempty"`      //
	Node      int           `json:"node,omitempty"`      //
	DB        *url.URL      `json:"db,omitempty"`        //
}

// runType represents operational settings
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
