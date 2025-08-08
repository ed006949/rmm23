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
		structMap map[string]mod_reflect.FieldTypeInfo
	)

	switch structMap, err = mod_reflect.BuildStructMap(out, TagName); {
	case err != nil:
		return
	}

	fmt.Print(structMap)

	var (
		interimEntry  = make(map[string][]string)
		outboundEntry = make(map[string]any)
		interim       []byte
	)

	interimEntry[mod_strings.F_dn.String()] = []string{e.DN}

	for _, b := range e.Attributes {
		var (
			structElement, ok = structMap[b.Name]
		)
		switch {
		case !ok:
			continue
		}

		switch {
		// Handle time.Time and []time.Time fields (LDAP GeneralizedTime)
		case (structElement.Kind == reflect.Slice &&
			structElement.ElemType == reflect.TypeOf(time.Time{})) ||
			(structElement.Type == reflect.TypeOf(time.Time{})):
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

	for a, b := range interimEntry {
		switch {
		case structMap[a].Kind == reflect.Slice:
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
