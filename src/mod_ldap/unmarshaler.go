package mod_ldap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"

	"rmm23/src/mod_reflect"
	"rmm23/src/mod_strings"
)

func UnmarshalEntry(e *ldap.Entry, out any) (err error) {
	var (
		structKind, structType = mod_reflect.BuildStructMap(out, "ldap")
	)

	fmt.Print(structKind, structType)

	var (
		interimEntry  = make(map[string][]string)
		outboundEntry = make(map[string]any)
		interim       []byte
	)

	for _, b := range e.Attributes {
		switch b.Name {
		case mod_strings.F_modifyTimestamp.String(), mod_strings.F_createTimestamp.String(): // LDAP <> time.Time workaround
			var (
				interimValues = make([]string, len(b.Values), len(b.Values))
			)
			for c, d := range b.Values {
				var (
					forErr      error
					interimTime time.Time
					interimData []byte
				)
				switch interimTime, forErr = ber.ParseGeneralizedTime([]byte(d)); {
				case forErr == nil:
					interimValues[c] = interimTime.Format(time.RFC3339)

					continue
				}

				interimValues[c] = string(interimData)
			}

			interimEntry[b.Name] = interimValues
		default:
			interimEntry[b.Name] = b.Values
		}
	}

	interimEntry[mod_strings.F_dn.String()] = []string{e.DN}

	for a, b := range interimEntry {
		switch {
		case structKind[a] == reflect.Slice && structType[a] != reflect.Struct:
			outboundEntry[a] = b
		default:
			switch {
			case len(b) > 0:
				outboundEntry[a] = b[0]
			}
		}
	}

	switch interim, err = json.Marshal(outboundEntry); {
	case err != nil:
		return
	}

	switch err = json.Unmarshal(interim, out); {
	case err != nil:
		return
	}

	return
}

// // UnmarshalText and all `attrTime` is for LDAP "specific" behavior
//
//	func (r *attrTime) UnmarshalText(inbound []byte) (err error) {
//		switch swInterim, swErr := ber.ParseGeneralizedTime(inbound); {
//		case swErr == nil:
//			r.Time = swInterim
//			return
//		}
//		var (
//			interim time.Time
//		)
//		switch err = interim.UnmarshalBinary(inbound); {
//		case err != nil:
//			return
//		}
//
//		r.Time = interim
//
//		return
//	}
//
// Function to walk struct fields and print kind
