package mod_db

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"rmm23/src/mod_slices"
)

func (r *AttrUUID) String() (outbound string)             { return r.UUID.String() }
func (r *AttrUUID) Entry() (outbound string)              { return entryKeyHeader + ":" + r.UUID.String() }
func (r *AttrUUID) Generate(space uuid.UUID, data []byte) { *r = AttrUUID{uuid.NewSHA1(space, data)} }

func (r *AttrUUID) MarshalJSON() (outbound []byte, err error) {
	return []byte(fmt.Sprintf("%q", r.String())), nil
}

func (r *AttrUUID) UnmarshalJSON(inbound []byte) (err error) {
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

	*r = AttrUUID{interimUUID}

	return
}

func (r *AttrUUID) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim uuid.UUID
		)
		switch interim, err = uuid.Parse(value); {
		case err != nil:
			continue
		}

		r.UUID = interim

		return // return only first value
	}

	return
}
