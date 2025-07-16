package mod_log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (r Z) MarshalZerologObject(e *zerolog.Event) {
	for a, b := range r {
		// switch a {
		// case E:
		// 	a = zerolog.ErrorFieldName
		// case M:
		// 	a = zerolog.MessageFieldName
		// 	// case T:
		// 	// 	a = zerolog.TypeFieldName
		// }

		switch value := b.(type) {
		// case nameValue:
		// e.Str(a, value.String())
		// case configValue:
		// e.Str(a, value.String())
		// case dryRunFlag:
		// e.Bool(a, value.Flag())
		// case modeValue:
		// e.Str(a, value.String())
		// case verbosityLevel:
		// e.Str(a, value.String())
		case error:
			e.AnErr(a, value)
		case []error:
			e.Errs(a, value)
		default:
			switch a {
			case T:
				e.Type(a, b)
			default:
				e.Interface(a, value)
			}
		}
	}

	switch {
	// case DryRun.Flag():
	// 	e.Bool(DryRun.Name(), DryRun.Flag())
	}
}

func (r Z) Emergency()     { log.Fatal().EmbedObject(r).Send() } // rfc3164 ----
func (r Z) Alert()         { log.Fatal().EmbedObject(r).Send() } // rfc3164 ----
func (r Z) Critical()      { log.Fatal().EmbedObject(r).Send() } // rfc3164 ----
func (r Z) Error()         { log.Error().EmbedObject(r).Send() } // rfc3164 +
func (r Z) Warning()       { log.Warn().EmbedObject(r).Send() }  // rfc3164 +
func (r Z) Notice()        { log.Info().EmbedObject(r).Send() }  // rfc3164 ----
func (r Z) Informational() { log.Info().EmbedObject(r).Send() }  // rfc3164 +
func (r Z) Debug()         { log.Debug().EmbedObject(r).Send() } // rfc3164 +
func (r Z) Trace()         { log.Trace().EmbedObject(r).Send() } // specific +
func (r Z) Panic()         { log.Panic().EmbedObject(r).Send() } // specific +
func (r Z) Quiet()         {}                                    // specific ----
func (r Z) Disabled()      {}                                    // specific ----
