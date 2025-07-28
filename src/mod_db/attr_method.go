package mod_db

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (r *attrDN) String() (outbound string) { return string(*r) }

// uuid.UUID

func (r *attrUUID) String() (outbound string) { return uuid.UUID(*r).String() }
func (r *attrUUID) Entry() (outbound string)  { return entryKeyHeader + ":" + r.String() }

// // MarshalText implements encoding.TextMarshaler
// func (r attrUUID) MarshalText() ([]byte, error) { return []byte(r.String()), nil }
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

// String implements fmt.Stringer.
//
// convert time.Time to String() seconds.
func (r *attrTime) String() (outbound string) { return strconv.FormatInt(time.Time(*r).Unix(), 10) }

// MarshalText implements encoding.TextMarshaler.
//
// convert time.Time to seconds.
func (r *attrTime) MarshalText() (outbound []byte, err error) { return []byte(r.String()), nil }

// UnmarshalText implements encoding.TextUnmarshaler.
//
// convert seconds to time.Time.
func (r *attrTime) UnmarshalText(inbound []byte) (err error) {
	var (
		interim int64
	)

	switch interim, err = strconv.ParseInt(string(inbound), 10, 64); {
	case err != nil:
		return
	}

	*r = attrTime(time.Unix(interim, int64(time.Second)))

	return
}
