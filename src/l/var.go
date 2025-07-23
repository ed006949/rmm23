package l

import (
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"rmm23/src/mod_bools"
	"rmm23/src/mod_errors"
)

var (
	buildName      string                       // to be set by builder
	buildVerbosity = zerolog.InfoLevel.String() // defaults
	buildDryRun    = "true"                     // defaults
	buildMode      = "0"                        // defaults
	buildNode      = "0"                        // defaults
	buildDB        = "redis://localhost:6379"   // defaults
	buildConfig    = buildName + ".json"        // defaults
	buildTime      string                       // to be set by builder
	buildCommit    string                       // to be set by builder
)
var (
	daemonFlagName = map[int]string{
		// daemonName:      "name",
		daemonVerbosity: "verbosity",
		daemonDryRun:    "dry-run",
		daemonMode:      "mode",
		daemonNode:      "node",
		daemonDB:        "db",
		daemonConfig:    "config",
		// daemonTime:      "time",
		// daemonCommit:    "commit",
	}
	daemonEnvName = map[int]string{
		// daemonName:      envName(daemonName),
		daemonVerbosity: envName(daemonVerbosity),
		daemonDryRun:    envName(daemonDryRun),
		daemonMode:      envName(daemonMode),
		daemonNode:      envName(daemonNode),
		daemonDB:        envName(daemonDB),
		daemonConfig:    envName(daemonConfig),
		// daemonTime:      envName(daemonTime),
		// daemonCommit:    envName(daemonCommit),
	}
	daemonEnvDescription = map[int]string{
		// daemonName:      "daemon name (" + daemonEnvName[daemonName] + "=\"" + daemonEnvDefined[daemonName] + "\")",
		daemonVerbosity: "verbosity level (" + daemonEnvName[daemonVerbosity] + "=\"" + daemonEnvDefined[daemonVerbosity] + "\")",
		daemonDryRun:    "dry-run flag (" + daemonEnvName[daemonDryRun] + "=\"" + daemonEnvDefined[daemonDryRun] + "\")",
		daemonMode:      "operational mode (" + daemonEnvName[daemonMode] + "=\"" + daemonEnvDefined[daemonMode] + "\")",
		daemonNode:      "node (" + daemonEnvName[daemonNode] + "=\"" + daemonEnvDefined[daemonNode] + "\")",
		daemonDB: "db url (" + daemonEnvName[daemonDB] + "=\"" + mod_errors.StripErr1(url.Parse(daemonEnvDefined[daemonDB])).Redacted() + "\")\n" +
			"\"redis://username:password@redis-host:6379\"\n" +
			"\"redis-sentinel://username:password@redis-sentinel-host1:6379,redis-sentinel-host2:6379\"",
		daemonConfig: "config file (" + daemonEnvName[daemonConfig] + "=\"" + daemonEnvDefined[daemonConfig] + "\")",
		// daemonTime:   "build time (" + daemonEnvName[daemonTime] + "=\"" + daemonEnvDefined[daemonTime] + "\")",
		// daemonCommit: "commit hash (" + daemonEnvName[daemonCommit] + "=\"" + daemonEnvDefined[daemonCommit] + "\")",
	}
	daemonEnvDefined = map[int]string{
		// daemonName:      os.Getenv(daemonEnvName[daemonName]),
		daemonVerbosity: os.Getenv(daemonEnvName[daemonVerbosity]),
		daemonDryRun:    os.Getenv(daemonEnvName[daemonDryRun]),
		daemonMode:      os.Getenv(daemonEnvName[daemonMode]),
		daemonNode:      os.Getenv(daemonEnvName[daemonNode]),
		daemonDB:        os.Getenv(daemonEnvName[daemonDB]),
		daemonConfig:    os.Getenv(daemonEnvName[daemonConfig]),
		// daemonTime:      os.Getenv(daemonEnvName[daemonTime]),
		// daemonCommit:    os.Getenv(daemonEnvName[daemonCommit]),
	}
)
var (
	Run = runType{
		name: buildName,
		verbosity: func() zerolog.Level {
			switch value := daemonEnvDefined[daemonVerbosity]; {
			case len(value) != 0:
				return mod_errors.PanicErr1(zerolog.ParseLevel(value))
			default:
				return mod_errors.PanicErr1(zerolog.ParseLevel(buildVerbosity))
			}
		}(),
		dryRun: func() bool {
			switch value := daemonEnvDefined[daemonDryRun]; {
			case len(value) != 0:
				return mod_errors.PanicErr1(mod_bools.Parse(value))
			default:
				return mod_errors.PanicErr1(mod_bools.Parse(buildDryRun))
			}
		}(),
		mode: func() int {
			switch value := daemonEnvDefined[daemonMode]; {
			case len(value) != 0:
				return mod_errors.PanicErr1(strconv.Atoi(value))
			default:
				return mod_errors.PanicErr1(strconv.Atoi(buildMode))
			}
		}(),
		node: func() int {
			switch value := daemonEnvDefined[daemonNode]; {
			case len(value) != 0:
				return mod_errors.PanicErr1(strconv.Atoi(value))
			default:
				return mod_errors.PanicErr1(strconv.Atoi(buildNode))
			}
		}(),
		db: func() *url.URL {
			switch value := daemonEnvDefined[daemonDB]; {
			case len(value) != 0:
				return mod_errors.PanicErr1(url.Parse(value))
			default:
				return mod_errors.PanicErr1(url.Parse(buildDB))
			}
		}(),
		config: buildConfig,
		time: func() time.Time {
			switch value := daemonEnvDefined[daemonTime]; {
			case len(value) != 0:
				return mod_errors.StripErr1(time.Parse(time.RFC3339, value))
			default:
				return mod_errors.StripErr1(time.Parse(time.RFC3339, buildTime))
			}
		}(),
		commit: buildCommit,
	}
)
