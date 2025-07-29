package mod_crypto

import (
	"encoding/json"
)

func (r *Certificate) MarshalJSON() (outbound []byte, err error) {
	var (
		pfxData []byte
	)

	switch pfxData, err = r.EncodeP12(); {
	case err != nil:
		return
	}

	return json.Marshal(pfxData)
}

func (r *Certificate) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim []byte
	)

	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	return r.DecodeP12(interim)
}
