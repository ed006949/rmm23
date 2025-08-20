package mod_time

import (
	"encoding/json/v2"
	"strconv"
	"time"

	ber "github.com/go-asn1-ber/asn1-ber"

	"rmm23/src/mod_errors"
)

func (r *Time) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim int64
	)
	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	*r = Time{time.Unix(interim, 0)}

	return
}

func (r *Time) MarshalJSON() (outbound []byte, err error) {
	var (
		interim = r.Time.Unix()
	)
	switch outbound, err = json.Marshal(&interim); {
	case err != nil:
		return
	}

	return
}

func (r *Time) UnmarshalText(inbound []byte) (err error) {
	var (
		t time.Time
		i int64
	)

	// // Empty value = set zero
	// switch {
	// case len(inbound) == 0:
	// 	r.Time = t
	//
	// 	return
	// }

	// Try string-encoded int64
	switch i, err = strconv.ParseInt(string(inbound), 10, 64); {
	case err == nil:
		r.Time = time.Unix(i, 0)

		return
	}

	// Try BER GeneralizedTime first (handles LDAP timestamps)
	switch t, err = ber.ParseGeneralizedTime(inbound); {
	case err == nil:
		r.Time = t

		return
	}

	t = time.Time{}

	// Try time.Time's UnmarshalText (handles RFC3339, etc.)
	switch err = t.UnmarshalText(inbound); {
	case err == nil:
		r.Time = t

		return
	}

	// Try time.Time's UnmarshalBinary
	switch err = t.UnmarshalBinary(inbound); {
	case err == nil:
		r.Time = t

		return
	}

	return mod_errors.EParse
}
