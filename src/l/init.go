package l

import (
	"flag"
	"net/url"
	"strconv"

	"github.com/rs/zerolog"

	"rmm23/src/mod_bools"
)

func init() {
	// l.log function call nesting depth is 1
	zerolog.CallerSkipFrameCount = zerolog.CallerSkipFrameCount + 1

	// flag.Func(daemonEnvName[daemonName], daemonEnvDescription[daemonName], func(inbound string) (err error) { return nil })

	flag.Func(daemonFlagName[daemonVerbosity], daemonEnvDescription[daemonVerbosity], func(inbound string) (err error) {
		var (
			interim zerolog.Level
		)
		switch interim, err = zerolog.ParseLevel(inbound); {
		case err != nil:
			return err
		}
		Run.verbosity = interim
		return
	})
	flag.Func(daemonFlagName[daemonDryRun], daemonEnvDescription[daemonDryRun], func(inbound string) (err error) {
		var (
			interim bool
		)
		switch interim, err = mod_bools.Parse(inbound); {
		case err != nil:
			return err
		}
		Run.dryRun = interim
		return
	})
	flag.Func(daemonFlagName[daemonMode], daemonEnvDescription[daemonMode], func(inbound string) (err error) {
		var (
			interim int
		)
		switch interim, err = strconv.Atoi(inbound); {
		case err != nil:
			return err
		}
		Run.mode = interim
		return
	})
	flag.Func(daemonFlagName[daemonNode], daemonEnvDescription[daemonNode], func(inbound string) (err error) {
		var (
			interim int
		)
		switch interim, err = strconv.Atoi(inbound); {
		case err != nil:
			return err
		}
		Run.node = interim
		return
	})
	flag.Func(daemonFlagName[daemonDB], daemonEnvDescription[daemonDB], func(inbound string) (err error) {
		var (
			interim *url.URL
		)
		switch interim, err = url.Parse(inbound); {
		case err != nil:
			return err
		}
		Run.db = interim
		return
	})
	flag.Func(daemonFlagName[daemonConfig], daemonEnvDescription[daemonConfig], func(inbound string) (err error) {
		Run.config = inbound
		return
	})

	flag.Parse()
}
