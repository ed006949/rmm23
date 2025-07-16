package mod_log

import (
	"github.com/rs/zerolog/log"
)

func Emergency(e Object)     { log.Fatal().EmbedObject(e).Send() }
func Alert(e Object)         { log.Fatal().EmbedObject(e).Send() }
func Critical(e Object)      { log.Fatal().EmbedObject(e).Send() }
func Error(e Object)         { log.Error().EmbedObject(e).Send() }
func Warning(e Object)       { log.Warn().EmbedObject(e).Send() }
func Notice(e Object)        { log.Info().EmbedObject(e).Send() }
func Informational(e Object) { log.Info().EmbedObject(e).Send() }
func Debug(e Object)         { log.Debug().EmbedObject(e).Send() }
func Trace(e Object)         { log.Trace().EmbedObject(e).Send() }
func Panic(e Object)         { log.Panic().EmbedObject(e).Send() }
func Quiet(e Object)         {}
func Disabled(e Object)      {}
