package mod_log

import (
	"github.com/rs/zerolog"
)

const (
	E = "error"   // zerolog.ErrorFieldName hook
	M = "message" // zerolog.MessageFieldName hook
	T = "type"    // zerolog.TypeFieldName hook
)

const (
	Emergency     = verbosityLevel(zerolog.FatalLevel)
	Alert         = verbosityLevel(zerolog.FatalLevel)
	Critical      = verbosityLevel(zerolog.FatalLevel)
	Error         = verbosityLevel(zerolog.ErrorLevel)
	Warning       = verbosityLevel(zerolog.WarnLevel)
	Notice        = verbosityLevel(zerolog.InfoLevel)
	Informational = verbosityLevel(zerolog.InfoLevel)
	Debug         = verbosityLevel(zerolog.DebugLevel)
	Trace         = verbosityLevel(zerolog.TraceLevel)
	Panic         = verbosityLevel(zerolog.PanicLevel)
	Quiet         = verbosityLevel(zerolog.NoLevel)
	Disabled      = verbosityLevel(zerolog.Disabled)
)
