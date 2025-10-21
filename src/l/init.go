package l

import (
	"flag"

	"github.com/rs/zerolog"
)

func Initialize() {
	zerolog.CallerSkipFrameCount = zerolog.CallerSkipFrameCount + 1 /* level method */ + 1 /* log method */

	Run.verbositySet(Run.verbosity)

	flag.Func(daemonFlagName[daemonConfig], daemonEnvDescription[daemonConfig], Run.configSetString)
	flag.Func(daemonFlagName[daemonDB], daemonEnvDescription[daemonDB], Run.dbSetString)
	flag.Func(daemonFlagName[daemonDryRun], daemonEnvDescription[daemonDryRun], Run.dryRunSetString)
	flag.Func(daemonFlagName[daemonMode], daemonEnvDescription[daemonMode], Run.modeSetString)
	flag.Func(daemonFlagName[daemonNode], daemonEnvDescription[daemonNode], Run.nodeSetString)
	flag.Func(daemonFlagName[daemonVerbosity], daemonEnvDescription[daemonVerbosity], Run.verbositySetString)

	flag.Parse()
}
