package mod_db

import (
	"strconv"
	"time"

	ber "github.com/go-asn1-ber/asn1-ber"
)

// // MarshalText implements encoding.TextMarshaler
// func (r *attrUUID) MarshalText() ([]byte, error) { return []byte(r.String()), nil }
//
// // UnmarshalText implements encoding.TextUnmarshaler
// func (r *attrUUID) UnmarshalText(inbound []byte) (err error) {
// 	var (
// 		interim uuid.UUID
// 	)
//
// 	switch interim, err = uuid.Parse(string(inbound)); {
// 	case err != nil:
// 		return
// 	}
//
// 	*r = attrUUID(interim)
// 	return
// }

// MarshalText implements encoding.TextMarshaler.
func (r *attrTime) MarshalText() (outbound []byte, err error) {
	return []byte(strconv.FormatInt(time.Time(*r).Unix(), 10)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (r *attrTime) UnmarshalText(inbound []byte) (err error) {
	var (
		interim time.Time
	)

	switch interim, err = ber.ParseGeneralizedTime(inbound); {
	case err != nil:
		return
	}

	*r = attrTime(interim)

	return
}
