package mod_db

import (
	"encoding/json"
	"fmt"
	"time"

	ber "github.com/go-asn1-ber/asn1-ber"

	"rmm23/src/mod_slices"
)

func (r *attrTime) String() (outbound string) { return r.Time.String() }

func (r *attrTime) MarshalJSON() (outbound []byte, err error) {
	return []byte(fmt.Sprintf("%d", r.Time.Unix())), nil
}

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

func (r *attrTime) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			interim time.Time
		)
		switch interim, err = ber.ParseGeneralizedTime([]byte(value)); {
		case err != nil:
			continue
		}

		r.Time = interim

		return // return only first value
	}

	return
}
