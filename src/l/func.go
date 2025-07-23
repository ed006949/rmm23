package l

import (
	"strings"
)

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

func envName(inbound int) (outbound string) {
	return strings.ReplaceAll(
		strings.ToUpper(
			string(buildName)+"_"+daemonFlagName[inbound]),
		"-",
		"_",
	)
}
