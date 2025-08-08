package mod_db

import (
	"encoding/json"
	"time"
)

type attrTime struct{ time.Time }

func (r *attrTime) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim int64
	)
	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	*r = attrTime{time.Unix(interim, 0)}

	return
}

func (r *attrTime) MarshalJSON() (outbound []byte, err error) {
	var (
		interim = r.Time.Unix()
	)
	switch outbound, err = json.Marshal(&interim); {
	case err != nil:
		return
	}

	return
}

// func (r *attrTime) UnmarshalText(inbound []byte) (err error) {
// 	var (
// 		t time.Time
// 		i int64
// 	)
//
// 	// // Empty value = set zero
// 	// switch {
// 	// case len(inbound) == 0:
// 	// 	*r = attrTime(t.Unix())
// 	//
// 	// 	return
// 	// }
//
// 	// Try string-encoded int64
// 	switch i, err = strconv.ParseInt(string(inbound), 10, 64); {
// 	case err == nil:
// 		*r = attrTime(i)
//
// 		return
// 	}
//
// 	// Try BER GeneralizedTime first (handles LDAP timestamps)
// 	switch t, err = ber.ParseGeneralizedTime(inbound); {
// 	case err == nil:
// 		*r = attrTime(t.Unix())
//
// 		return
// 	}
//
// 	t = time.Time{}
//
// 	// Try time.Time's built-in UnmarshalText (handles RFC3339, etc.)
// 	switch err = t.UnmarshalText(inbound); {
// 	case err == nil:
// 		*r = attrTime(t.Unix())
//
// 		return
// 	}
//
// 	// Try time.Time's built-in UnmarshalBinary
// 	switch err = t.UnmarshalBinary(inbound); {
// 	case err == nil:
// 		*r = attrTime(t.Unix())
//
// 		return
// 	}
//
// 	return mod_errors.EParse
// }
//
// // func (r *attrTime) MarshalText() (outbound []byte, err error) {
// // 	return fmt.Appendf(nil, "%d", *r), nil
// // }
//
// func (r *attrTime) Time() (outbound time.Time) { return time.Unix(int64(*r), 0) }
//
// func (r *attrTime) Set(inbound time.Time) { *r = attrTime(inbound.Unix()) }
