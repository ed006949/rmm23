package l

import (
	"flag"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"rmm23/src/mod_bools"
	"rmm23/src/mod_errors"
)

func init() {
	// l.log function call nesting depth is 1
	zerolog.CallerSkipFrameCount = zerolog.CallerSkipFrameCount + 1

	// 	daemonName = iota
	//	daemonVerbosity
	//	daemonDryRun
	//	daemonMode
	//	daemonNode
	//	daemonDB
	//	daemonConfig
	//	daemonTime
	//	daemonCommit
	var (
		// cliName      = flag.String(
		// 	daemonParamDescription[daemonName],
		// 	os.Getenv(daemonEnvName[daemonName]),
		// 	daemonEnvDescription[daemonName],
		// )
		cliVerbosity = flag.String(
			daemonParamDescription[daemonVerbosity],
			mod_errors.StripErr1(zerolog.ParseLevel(os.Getenv(daemonEnvName[daemonVerbosity]))).String(),
			daemonEnvDescription[daemonVerbosity],
		)
		cliDryRun = flag.Bool(
			daemonParamDescription[daemonDryRun],
			mod_errors.StripErr1(mod_bools.Parse(os.Getenv(daemonEnvName[daemonDryRun]))),
			daemonEnvDescription[daemonDryRun],
		)
		cliMode = flag.Int(
			daemonParamDescription[daemonMode],
			mod_errors.StripErr1(strconv.Atoi(os.Getenv(daemonEnvName[daemonMode]))),
			daemonEnvDescription[daemonMode],
		)
		cliNode = flag.Int(
			daemonParamDescription[daemonNode],
			mod_errors.StripErr1(strconv.Atoi(os.Getenv(daemonEnvName[daemonNode]))),
			daemonEnvDescription[daemonNode],
		)
		cliDB = flag.String(
			daemonParamDescription[daemonDB],
			os.Getenv(daemonEnvName[daemonDB]),
			daemonEnvDescription[daemonDB],
		)
		cliConfig = flag.String(
			daemonParamDescription[daemonConfig],
			os.Getenv(daemonEnvName[daemonConfig]),
			daemonEnvDescription[daemonConfig],
		)
	)
	flag.Parse()

	switch FlagIsFlagExist(daemonEnvName[daemonVerbosity]); {
	case true:
	case false:
	}
	switch FlagIsFlagExist(daemonEnvName[daemonDryRun]); {
	case true:
		run.dryRun = *cliDryRun
	case false:
		run.dryRun = mod_errors.StripErr1(mod_bools.Parse(os.Getenv(daemonEnvName[daemonDryRun])))
	}
	switch FlagIsFlagExist(daemonEnvName[daemonMode]); {
	case true:
	case false:
	}
	switch FlagIsFlagExist(daemonEnvName[daemonNode]); {
	case true:
	case false:
	}
	switch FlagIsFlagExist(daemonEnvName[daemonDB]); {
	case true:
	case false:
	}
	switch FlagIsFlagExist(daemonEnvName[daemonConfig]); {
	case true:
	case false:
	}
}
func InitCLI() {
	var (
		cliName      = flag.String(Name.Name(), os.Getenv(Name.EnvName()), Name.EnvDescription())
		cliConfig    = flag.String(Config.Name(), os.Getenv(Config.EnvName()), Config.EnvDescription())
		cliDryRun    = flag.Bool(DryRun.Name(), mod_errors.StripErr1(mod_bools.Parse(os.Getenv(DryRun.EnvName()))), DryRun.EnvDescription())
		cliMode      = flag.String(Mode.Name(), os.Getenv(Mode.EnvName()), Mode.EnvDescription())
		cliVerbosity = flag.String(Verbosity.Name(), os.Getenv(Verbosity.EnvName()), Verbosity.EnvDescription())
	)
	flag.Parse()

	switch {
	case FlagIsFlagExist(Name.Name()):
		Name.Set(*cliName)
		fallthrough
	case FlagIsFlagExist(Config.Name()):
		Config.Set(*cliConfig)
		fallthrough
	case FlagIsFlagExist(DryRun.Name()):
		DryRun.Set(*cliDryRun)
		fallthrough
	case FlagIsFlagExist(Mode.Name()):
		Mode.Set(*cliMode)
		fallthrough
	case FlagIsFlagExist(Verbosity.Name()):
		Verbosity.Set(*cliVerbosity)
		fallthrough
	default:
	}
}
