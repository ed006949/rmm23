package l

import (
	"encoding/json/v2"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (r *runType) verbositySet(inbound zerolog.Level) {
	r.verbosity = inbound
	log.Logger = log.Level(r.verbosity).With().Timestamp().Caller().Ctx(r.ctx).Logger().Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    false,
		TimeFormat: time.RFC3339,
		// FormatFieldValue: func(i any) string { return fmt.Sprintf("\"%s\"", i) },
	})
}

func (r *runType) configSetString(inbound string) (err error) {
	Run.config = inbound

	return
}
func (r *runType) dbSetString(inbound string) (err error) {
	var (
		interim *url.URL
	)
	switch interim, err = url.Parse(inbound); {
	case err != nil:
		log.Error().Err(err).Send()

		return
	}

	Run.db = interim

	return
}
func (r *runType) dryRunSetString(inbound string) (err error) {
	var (
		interim bool
	)
	switch interim, err = strconv.ParseBool(inbound); {
	case err != nil:
		log.Error().Err(err).Send()

		return
	}

	r.dryRun = interim

	return
}
func (r *runType) modeSetString(inbound string) (err error) {
	var (
		interim int
	)
	switch interim, err = strconv.Atoi(inbound); {
	case err != nil:
		log.Error().Err(err).Send()

		return
	}

	Run.mode = interim

	return
}
func (r *runType) nodeSetString(inbound string) (err error) {
	var (
		interim int
	)
	switch interim, err = strconv.Atoi(inbound); {
	case err != nil:
		log.Error().Err(err).Send()

		return
	}

	Run.node = interim

	return
}
func (r *runType) verbositySetString(inbound string) (err error) {
	var (
		interim zerolog.Level
	)
	switch interim, err = zerolog.ParseLevel(inbound); {
	case err != nil:
		log.Error().Err(err).Send()

		return
	}

	Run.verbositySet(interim)

	return
}

func (r *runType) ConfigUnmarshal(inbound any) (err error) {
	var (
		content []byte
	)
	switch content, err = os.ReadFile(r.config); {
	case err != nil:
		log.Error().Err(err).Send()

		return
	}

	// switch err = json.Unmarshal(content, inbound); {
	// case err != nil:
	// 	return
	// }
	//
	// return

	return json.Unmarshal(content, inbound)
}

func (r *runType) BuildTimeValue() (outbound string)  { return r.time.String() }
func (r *runType) CommitHashValue() (outbound string) { return r.commit }
func (r *runType) DryRunValue() (outbound bool)       { return r.dryRun }
func (r *runType) NameValue() (outbound string)       { return r.name }

func (r *runType) DryRunName() (outbound string) { return daemonFlagName[daemonDryRun] }
