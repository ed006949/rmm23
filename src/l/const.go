package l

import (
	"github.com/rs/zerolog"
)

const (
	E = "error"   // zerolog.ErrorFieldName hook
	M = "message" // zerolog.MessageFieldName hook
	T = "type"    // zerolog.TypeFieldName hook
)

const (
	Name      nameType      = "name"
	Config    configType    = "config"
	DryRun    dryRunType    = "dry-run"
	Mode      modeType      = "mode"
	Verbosity verbosityType = "verbosity"
	GitCommit gitCommitType = "buildCommit"
)
const (
	NoDryRun dryRunValue = false
	DoDryRun dryRunValue = true
)
const (
	Init modeValue = iota
	Deploy
	CLI
	Daemon
)
const (
	Emergency     = verbosityValue(zerolog.FatalLevel)
	Alert         = verbosityValue(zerolog.FatalLevel)
	Critical      = verbosityValue(zerolog.FatalLevel)
	Error         = verbosityValue(zerolog.ErrorLevel)
	Warning       = verbosityValue(zerolog.WarnLevel)
	Notice        = verbosityValue(zerolog.InfoLevel)
	Informational = verbosityValue(zerolog.InfoLevel)
	Debug         = verbosityValue(zerolog.DebugLevel)
	Trace         = verbosityValue(zerolog.TraceLevel)
	Panic         = verbosityValue(zerolog.PanicLevel)
	Quiet         = verbosityValue(zerolog.NoLevel)
	Disabled      = verbosityValue(zerolog.Disabled)
)
const (
	daemonName = iota
	daemonVerbosity
	daemonDryRun
	daemonMode
	daemonNode
	daemonDB
	daemonConfig
	daemonTime
	daemonCommit
)
