package mod_db

import (
	"strconv"
	"time"

	ber "github.com/go-asn1-ber/asn1-ber"

	"rmm23/src/mod_errors"
)

type attrTime struct{ time.Time }

func (r *attrTime) UnmarshalText(inbound []byte) (err error) {
	var (
		interim time.Time
		i       int64
	)

	// Empty value = set zero
	switch {
	case len(inbound) == 0:
		r.Time = interim

		return
	}

	// Try BER GeneralizedTime first (handles LDAP timestamps)
	switch interim, err = ber.ParseGeneralizedTime(inbound); {
	case err == nil:
		r.Time = interim

		return
	}

	interim = time.Time{}

	// Try time.Time's built-in UnmarshalText (handles RFC3339, etc.)
	switch err = interim.UnmarshalText(inbound); {
	case err == nil:
		r.Time = interim

		return
	}

	// Try time.Time's built-in UnmarshalBinary
	switch err = interim.UnmarshalBinary(inbound); {
	case err == nil:
		r.Time = interim

		return
	}

	// Try string-encoded int64
	switch i, err = strconv.ParseInt(string(inbound), 10, 64); {
	case err == nil:
		interim = time.Unix(i, 0)
		r.Time = interim

		return
	}

	return mod_errors.EParse
}

func (r *attrTime) MarshalText() (outbound []byte, err error) {
	return []byte(strconv.FormatInt(r.Time.Unix(), 10)), nil
}

// func (r *attrTime) String() (outbound string) { return r.Time.String() }
