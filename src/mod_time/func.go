package mod_time

import (
	"github.com/rs/zerolog/log"
)

func UnmarshalText(inbound []byte) (outbound *Time, err error) {
	var (
		interim = new(Time)
	)
	switch err = interim.UnmarshalText(inbound); {
	case err != nil:
		log.Error().Err(err).Send()

		return
	}

	return interim, nil
}
