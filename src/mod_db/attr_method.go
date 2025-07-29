package mod_db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

//
// attrDN

func (r *attrDN) String() (outbound string) { return string(*r) }

//
// attrUUID

func (r *attrUUID) String() (outbound string) { return uuid.UUID(*r).String() }

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

	return nil
}

//
// attrIDNumber

func (r *attrIDNumber) MarshalJSON() (outbound []byte, err error) {
	return []byte(fmt.Sprintf("%d", *r)), nil
}
func (r *attrIDNumber) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim uint
	)

	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	*r = attrIDNumber(interim)

	return
}
