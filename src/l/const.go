package l

const (
	E = "error"   // zerolog.ErrorFieldName hook
	M = "message" // zerolog.MessageFieldName hook
	T = "type"    // zerolog.TypeFieldName hook
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

// const (
//	Emergency     = verbosityValue(zerolog.FatalLevel)
//	Alert         = verbosityValue(zerolog.FatalLevel)
//	Critical      = verbosityValue(zerolog.FatalLevel)
//	Error         = verbosityValue(zerolog.ErrorLevel)
//	Warning       = verbosityValue(zerolog.WarnLevel)
//	Notice        = verbosityValue(zerolog.InfoLevel)
//	Informational = verbosityValue(zerolog.InfoLevel)
//	Debug         = verbosityValue(zerolog.DebugLevel)
//	Trace         = verbosityValue(zerolog.TraceLevel)
//	Panic         = verbosityValue(zerolog.PanicLevel)
//	Quiet         = verbosityValue(zerolog.NoLevel)
//	Disabled      = verbosityValue(zerolog.Disabled)
// )

const (
	BulkOpsSize = 16
)
