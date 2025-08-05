package mod_bools

import (
	"encoding/json"
	"strings"
)

func (r *AttrBool) MarshalJSON() (outbound []byte, err error) { return json.Marshal(r.Int()) }

func (r *AttrBool) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim int
	)
	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	r.Parse(interim)

	return
}

func (r *AttrBool) Int() (outbound int) {
	switch *r {
	case true:
		return 1
	default:
		return 0
	}
}

func (r *AttrBool) Parse(inbound any) {
	switch v := inbound.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		switch v {
		case 0:
			*r = false
		case 1:
			*r = true
		}
	case string:
		switch strings.ToLower(v) {
		case "0", "f", "n", "false", "no", "off":
			*r = false
		case "1", "t", "y", "true", "yes", "on":
			*r = true
		}
	}
}
