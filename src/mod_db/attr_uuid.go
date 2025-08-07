package mod_db

import (
	"encoding/json"

	"github.com/google/uuid"

	"rmm23/src/mod_slices"
)

func (r *attrUUID) String() (outbound string) { return r.UUID.String() }
func (r *attrUUID) Entry() (outbound string) {
	return entryKeyHeader + headerSeparator + r.UUID.String()
}

func (r *attrUUID) generate(space uuid.UUID, data []byte) {
	*r = attrUUID{uuid.NewSHA1(space, data)}
}

func generateUUID(space uuid.UUID, data []byte) (outbound *attrUUID) {
	outbound = new(attrUUID)
	outbound.generate(space, data)

	return
}

func (r *attrUUID) MarshalJSON() (outbound []byte, err error) {
	return json.Marshal(r.UUID)
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

	*r = attrUUID{interimUUID}

	return
}

func (r *attrUUID) UnmarshalLDAPAttr(values []string) (err error) {
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
