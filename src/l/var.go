package l

import (
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"rmm23/src/mod_bools"
	"rmm23/src/mod_errors"
)

var (
	buildName      string                       // to be set by builder
	buildVerbosity = zerolog.InfoLevel.String() // defaults
	buildDryRun    = DoDryRun.String()          // defaults
	buildMode      = Init.String()              // defaults
	buildNode      = "0"                        // defaults
	buildDB        = "redis://localhost:6379"   // defaults
	buildConfig    = buildName + ".json"        // defaults
	buildTime      string                       // to be set by builder
	buildCommit    string                       // to be set by builder
)

var (
	control = &DaemonConfig{
		Name:      "",
		Verbosity: 0,
		DryRun:    DoDryRun,
		Mode:      0,
		Node:      0,
		DB:        &url.URL{},
		Config:    "",
		build: buildType{
			name:      buildName,
			verbosity: buildVerbosity,
			dryRun:    buildDryRun,
			mode:      buildMode,
			node:      buildNode,
			db:        buildDB,
			config:    buildConfig,
			time:      buildTime,
			commit:    buildCommit,
		},
	}
	run = runType{
		name:      buildName,
		verbosity: mod_errors.PanicErr1(zerolog.ParseLevel(buildVerbosity)),
		dryRun:    mod_errors.PanicErr1(mod_bools.Parse(buildDryRun)),
		mode:      mod_errors.PanicErr1(strconv.Atoi(buildMode)),
		node:      mod_errors.PanicErr1(strconv.Atoi(buildNode)),
		db:        mod_errors.PanicErr1(url.Parse(buildDB)),
		config:    buildConfig,
		time:      mod_errors.StripErr1(time.Parse(time.RFC3339, buildTime)),
		commit:    buildCommit,
	}
)
var (
	daemonParamDescription = map[int]string{
		daemonName:      "name",
		daemonVerbosity: "verbosity",
		daemonDryRun:    "dry-run",
		daemonMode:      "mode",
		daemonNode:      "node",
		daemonDB:        "db",
		daemonConfig:    "config",
		daemonTime:      "time",
		daemonCommit:    "commit",
	}
	daemonEnvName = map[int]string{
		daemonName:      envName(daemonName),
		daemonVerbosity: envName(daemonVerbosity),
		daemonDryRun:    envName(daemonDryRun),
		daemonMode:      envName(daemonMode),
		daemonNode:      envName(daemonNode),
		daemonDB:        envName(daemonDB),
		daemonConfig:    envName(daemonConfig),
		daemonTime:      envName(daemonTime),
		daemonCommit:    envName(daemonCommit),
	}
	daemonEnvDescription = map[int]string{
		daemonName:      "daemon name (" + daemonEnvName[daemonName] + ")",
		daemonVerbosity: "verbosity level (" + daemonEnvName[daemonVerbosity] + ")",
		daemonDryRun:    "dry-run flag (" + daemonEnvName[daemonDryRun] + ")",
		daemonMode:      "operational mode (" + daemonEnvName[daemonMode] + ")",
		daemonNode:      "node (" + daemonEnvName[daemonNode] + ")",
		daemonDB:        "db url (" + daemonEnvName[daemonDB] + ")",
		daemonConfig:    "config file (" + daemonEnvName[daemonConfig] + ")",
		daemonTime:      "build time (" + daemonEnvName[daemonTime] + ")",
		daemonCommit:    "commit hash (" + daemonEnvName[daemonCommit] + ")",
	}
)
var (
	dryRunDescription = map[dryRunValue]string{
		NoDryRun: "false",
		DoDryRun: "true",
	}
)
var (
	modeDescription = map[modeValue]string{
		Init:   "init",
		Deploy: "deploy",
		CLI:    "cli",
		Daemon: "daemon",
	}
)
