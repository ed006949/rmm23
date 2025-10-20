package l

import (
	"strings"

	"github.com/rs/zerolog/log"
)

func Emergency(e Z)     { log.Fatal().EmbedObject(e).Send() }
func Alert(e Z)         { log.Fatal().EmbedObject(e).Send() }
func Critical(e Z)      { log.Fatal().EmbedObject(e).Send() }
func Error(e Z)         { log.Error().EmbedObject(e).Send() }
func Warning(e Z)       { log.Warn().EmbedObject(e).Send() }
func Notice(e Z)        { log.Info().EmbedObject(e).Send() }
func Informational(e Z) { log.Info().EmbedObject(e).Send() }
func Debug(e Z)         { log.Debug().EmbedObject(e).Send() }
func Trace(e Z)         { log.Trace().EmbedObject(e).Send() }
func Panic(e Z)         { log.Panic().EmbedObject(e).Send() }
func Quiet(e Z)         {}
func Disabled(e Z)      {}

// envName generates an environment variable name based on the daemon name and a daemon flag name.
// It converts the combined name to uppercase and replaces hyphens with underscores.
func envName(inbound int) (outbound string) {
	return strings.ReplaceAll(
		strings.ToUpper(
			buildName+"_"+daemonFlagName[inbound]),
		"-",
		"_",
	)
}
