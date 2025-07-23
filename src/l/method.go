package l

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (r Z) MarshalZerologObject(e *zerolog.Event) {
	for a, b := range r {
		switch value := b.(type) {
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
	case Run.dryRun:
		e.Bool(daemonFlagName[daemonDryRun], Run.dryRun)
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

func (r *runType) verbositySet(inbound zerolog.Level) {
	r.verbosity = inbound
	log.Logger = log.Level(r.verbosity).With().Timestamp().Caller().Logger().Output(zerolog.ConsoleWriter{
		Out:              os.Stderr,
		NoColor:          false,
		TimeFormat:       time.RFC3339,
		FormatFieldValue: func(i interface{}) string { return fmt.Sprintf("\"%s\"", i) },
	})
}
func (r *runType) ConfigUnmarshal(inbound any) (err error) {
	var (
		content []byte
	)

	switch content, err = os.ReadFile(r.config); {
	case err != nil:
		return
	}

	return json.Unmarshal(content, inbound)
}

func (r *runType) DryRunValue() (outbound bool)  { return r.dryRun }
func (r *runType) NameValue() (outbound string)       { return r.name }
func (r *runType) CommitHashValue() (outbound string) { return r.commit }
func (r *runType) BuildTimeValue() (outbound string)  { return r.time.String() }
