package l

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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
		case nameValue:
			e.Str(a, value.String())
		case configValue:
			e.Str(a, value.String())
		case dryRunValue:
			e.Bool(a, value.Flag())
		case modeValue:
			e.Str(a, value.String())
		case verbosityValue:
			e.Str(a, value.String())
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
	case DryRun.Flag():
		e.Bool(DryRun.Name(), DryRun.Flag())
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

func (r nameValue) Set()      { control.set(Name, r) }      // Package predefined Flag hook
func (r configValue) Set()    { control.set(Config, r) }    // Package predefined Flag hook
func (r dryRunValue) Set()    { control.set(DryRun, r) }    // Package predefined Flag hook
func (r modeValue) Set()      { control.set(Mode, r) }      // Package predefined Flag hook
func (r verbosityValue) Set() { control.set(Verbosity, r) } // Package predefined Flag hook

func (nameType) Set(inbound any)      { control.set(Name, inbound) }      // Package Flag hook
func (configType) Set(inbound any)    { control.set(Config, inbound) }    // Package Flag hook
func (dryRunType) Set(inbound any)    { control.set(DryRun, inbound) }    // Package Flag hook
func (modeType) Set(inbound any)      { control.set(Mode, inbound) }      // Package Flag hook
func (verbosityType) Set(inbound any) { control.set(Verbosity, inbound) } // Package Flag hook

func (nameType) Name() string      { return string(Name) }      // Package Flag Name
func (configType) Name() string    { return string(Config) }    // Package Flag Name
func (dryRunType) Name() string    { return string(DryRun) }    // Package Flag Name
func (modeType) Name() string      { return string(Mode) }      // Package Flag Name
func (verbosityType) Name() string { return string(Verbosity) } // Package Flag Name
func (gitCommitType) Name() string { return string(GitCommit) } // Package Flag Name

func (nameType) Value() nameValue           { return control.Name }      // Package Flag Value
func (configType) Value() configValue       { return control.Config }    // Package Flag Value
func (dryRunType) Value() dryRunValue       { return control.DryRun }    // Package Flag Value
func (modeType) Value() modeValue           { return control.Mode }      // Package Flag Value
func (verbosityType) Value() verbosityValue { return control.Verbosity } // Package Flag Value
func (gitCommitType) Value() commitValue    { return control.Commit }    // Package Flag Value

func (r nameType) EnvName() string      { return envName(r.Name()) } // Package Flag Env Name
func (r configType) EnvName() string    { return envName(r.Name()) } // Package Flag Env Name
func (r dryRunType) EnvName() string    { return envName(r.Name()) } // Package Flag Env Name
func (r modeType) EnvName() string      { return envName(r.Name()) } // Package Flag Env Name
func (r verbosityType) EnvName() string { return envName(r.Name()) } // Package Flag Env Name

func (r nameType) EnvDescription() string      { return envDescription(r) }
func (r configType) EnvDescription() string    { return envDescription(r) }
func (r dryRunType) EnvDescription() string    { return envDescription(r) }
func (r modeType) EnvDescription() string      { return envDescription(r) }
func (r verbosityType) EnvDescription() string { return envDescription(r) }

func (r dryRunType) Flag() bool              { return r.Value().Flag() }   // Package Flag Flag Value
func (r verbosityType) Level() zerolog.Level { return r.Value().Level() }  // Package Flag Level Value
func (r nameType) String() string            { return r.Value().String() } // Package Flag String Value
func (r configType) String() string          { return r.Value().String() } // Package Flag String Value
func (r dryRunType) String() string          { return r.Value().String() } // Package Flag String Value
func (r modeType) String() string            { return r.Value().String() } // Package Flag String Value
func (r verbosityType) String() string       { return r.Value().String() } // Package Flag String Value
func (r gitCommitType) String() string       { return r.Value().String() } // Package Flag String Value

func (r dryRunValue) Flag() bool              { return bool(r) }                   // Package Flag flag
func (r verbosityValue) Level() zerolog.Level { return zerolog.Level(r) }          // Package Flag level
func (r nameValue) String() string            { return string(r) }                 // Package Flag description
func (r configValue) String() string          { return string(r) }                 // Package Flag description
func (r dryRunValue) String() string          { return dryRunDescription[r] }      // Package Flag description
func (r modeValue) String() string            { return modeDescription[r] }        // Package Flag description
func (r verbosityValue) String() string       { return zerolog.Level(r).String() } // Package Flag description
func (r commitValue) String() string          { return string(r) }                 // Package Flag description

func (r *DaemonConfig) set(inboundKey any, inboundValue any) {
	switch inboundKey.(type) {
	case nameType:
		switch value := inboundValue.(type) {
		case nameValue:
			r.Name = value
		case string:
			switch {
			case len(value) == 0:
				return
			}
			r.Name = nameValue(value)
		}

	case configType:
		switch value := inboundValue.(type) {
		case configValue:
			r.Config = value
		case string:
			switch {
			case len(value) == 0:
				return
			}
			r.Config = configValue(value)
		}

	case dryRunType:
		switch value := inboundValue.(type) {
		case dryRunValue:
			r.DryRun = value
		case bool:
			r.DryRun = dryRunValue(value)

		case string:
			switch {
			case len(value) == 0:
				return
			}
			value = strings.ToLower(value)
			switch value {
			case "1", "t", "y", "true", "yes", "on":
				r.DryRun = true
			case "0", "f", "n", "false", "no", "off":
				r.DryRun = false
			}
		}

	case modeType:
		switch value := inboundValue.(type) {
		case modeValue:
			r.Mode = value
		case int:
			r.Mode = modeValue(value)

		case string:
			switch {
			case len(value) == 0:
				return
			}
			value = strings.ToLower(value)
			for a, b := range modeDescription {
				switch {
				case value == b:
					r.Mode = a
					return
				}
			}
		}

	case verbosityType:
		switch value := inboundValue.(type) {
		case verbosityValue:
			r.Verbosity = value
		case int8:
			r.Verbosity = verbosityValue(value)

		case zerolog.Level:
			r.Verbosity = verbosityValue(value)

		case string:
			switch {
			case len(value) == 0:
				return
			}
			value = strings.ToLower(value)
			switch interim, err := zerolog.ParseLevel(value); {
			case err != nil:
				return
			default:
				r.Verbosity = verbosityValue(interim)
			}
		}

		zerolog.SetGlobalLevel(r.Verbosity.Level()) // .!. how it works ....
		log.Logger = log.Level(r.Verbosity.Level()).With().Timestamp().Caller().Logger().Output(zerolog.ConsoleWriter{
			Out:              os.Stderr,
			NoColor:          false,
			TimeFormat:       time.RFC3339,
			FormatFieldValue: func(i interface{}) string { return fmt.Sprintf("\"%s\"", i) },
		})

	}
}

func (r *nameValue) UnmarshalJSON(inbound []byte) (err error)      { return json.Unmarshal(inbound, r) }
func (r *configValue) UnmarshalJSON(inbound []byte) (err error)    { return json.Unmarshal(inbound, r) }
func (r *dryRunValue) UnmarshalJSON(inbound []byte) (err error)    { return json.Unmarshal(inbound, r) }
func (r *modeValue) UnmarshalJSON(inbound []byte) (err error)      { return json.Unmarshal(inbound, r) }
func (r *verbosityValue) UnmarshalJSON(inbound []byte) (err error) { return json.Unmarshal(inbound, r) }
