package l

import (
	"flag"

	"github.com/rs/zerolog"
)

func Initialize() {
	// l.log function call nesting depth is 1
	zerolog.CallerSkipFrameCount = zerolog.CallerSkipFrameCount + 1

	Run.verbositySet(Run.verbosity)

	flag.Func(daemonFlagName[daemonVerbosity], daemonEnvDescription[daemonVerbosity], Run.verbositySetString)
	flag.Func(daemonFlagName[daemonDryRun], daemonEnvDescription[daemonDryRun], Run.dryRunSetString)
	flag.Func(daemonFlagName[daemonMode], daemonEnvDescription[daemonMode], Run.modeSetString)
	flag.Func(daemonFlagName[daemonNode], daemonEnvDescription[daemonNode], Run.nodeSetString)
	flag.Func(daemonFlagName[daemonDB], daemonEnvDescription[daemonDB], Run.dbSetString)
	flag.Func(daemonFlagName[daemonConfig], daemonEnvDescription[daemonConfig], Run.configSetString)

	flag.Parse()
}
