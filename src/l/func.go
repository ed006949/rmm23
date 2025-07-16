package l

import (
	"flag"
	"strings"
)

// TODO move projects to []error reporting style

// func Emergency(e Z)     { log.Fatal().EmbedObject(e).Send() }
// func Alert(e Z)         { log.Fatal().EmbedObject(e).Send() }
// func Critical(e Z)      { log.Fatal().EmbedObject(e).Send() }
// func Error(e Z)         { log.Error().EmbedObject(e).Send() }
// func Warning(e Z)       { log.Warn().EmbedObject(e).Send() }
// func Notice(e Z)        { log.Info().EmbedObject(e).Send() }
// func Informational(e Z) { log.Info().EmbedObject(e).Send() }
// func Debug(e Z)         { log.Debug().EmbedObject(e).Send() }
// func Trace(e Z)         { log.Trace().EmbedObject(e).Send() }
// func Panic(e Z)         { log.Panic().EmbedObject(e).Send() }
// func Quiet(e Z)         {}
// func Disabled(e Z)      {}

func envDescription(inbound any) (outbound string) {
	switch value := inbound.(type) {
	case nameType:
		return "daemon name (" + value.EnvName() + ")"
	case configType:
		return "config file (" + value.EnvName() + ")"
	case dryRunType:
		return "dry-run flag, overrides config (" + value.EnvName() + ")"
	case verbosityType:
		return "verbosity level, overrides config (" + value.EnvName() + ")"
	case modeType:
		return "operational mode (" + value.EnvName() + ")"
	default:
		return
	}
}

func EnvName(inbound string) (outbound string) {
	return strings.ReplaceAll(
		strings.ToUpper(
			Name.String()+"_"+inbound),
		"-",
		"_",
	)
}

func FlagIsFlagExist(name string) (outbound bool) {
	flag.Visit(func(fn *flag.Flag) {
		switch {
		case fn.Name == name:
			outbound = true
		}
	})
	return
}
