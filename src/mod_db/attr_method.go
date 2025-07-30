package mod_db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"rmm23/src/mod_errors"
)

//
// attrEntryType

func (r *attrEntryType) String() (outbound string) { return entryTypeName[*r] }
func (r *attrEntryType) Parse(inbound string) (err error) {
	switch value, ok := entryTypeID[inbound]; {
	case !ok:
		return mod_errors.EUnknownType
	default:
		*r = value

		return
	}
}

//
// attrDN

func (r *attrDN) String() (outbound string) { return string(*r) }

//
// attrUUID
// store/retrieve uuid.UUID as []byte

func (r *attrUUID) String() (outbound string) { return uuid.UUID(*r).String() }
func (r *attrUUID) Entry() (outbound string)  { return entryKeyHeader + ":" + uuid.UUID(*r).String() }

func (r *attrUUID) MarshalJSON() (outbound []byte, err error) {
	return []byte(fmt.Sprintf("%q", r.String())), nil
}

func (r *attrUUID) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim string
	)

	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	var (
		interimUUID uuid.UUID
	)
	switch interimUUID, err = uuid.Parse(interim); {
	case err != nil:
		return
	}

	*r = attrUUID(interimUUID)

	return
}

//
// attrTime
// store/retrieve time.Time as int64 to utilize redisearch NUMERIC search [min max]

func (r *attrTime) String() (outbound string) { return time.Time(*r).String() }

func (r *attrTime) MarshalJSON() (outbound []byte, err error) {
	return []byte(fmt.Sprintf("%d", time.Time(*r).Unix())), nil
}

func (r *attrTime) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim int64
	)

	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	*r = attrTime(time.Unix(interim, 0))

	return
}
