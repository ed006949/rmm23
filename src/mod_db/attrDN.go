package mod_db

import (
	"strings"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
)

const (
	dnSeparator     = "="
	dnPathSeparator = ","
)

type attrDN []attrDNFV

type attrDNFV struct{ Field, Value string }

// func (r *attrDN) MarshalJSON() (outbound []byte, err error) { return r.Time.MarshalJSON() }
//
// func (r *attrDN) UnmarshalJSON(inbound []byte) (err error) {
// 	switch swInterim, swErr := ber.ParseGeneralizedTime(inbound); {
// 	case swErr == nil:
// 		r.Time = swInterim
// 		return
// 	}
// 	var (
// 		interim []time.Time
// 	)
// 	switch err = json.Unmarshal(inbound, &interim); {
// 	case err != nil:
// 		return
// 	}
//
// 	r.Time = interim[0]
//
// 	return
// }

func (r *attrDN) UnmarshalText(inbound []byte) (err error) {
	var (
		interimFVs = mod_strings.Split(string(inbound), dnPathSeparator, mod_slices.FlagFilterEmpty|mod_slices.FlagTrimSpace)
		interim    = make(attrDN, len(interimFVs), len(interimFVs))
	)
	for a, b := range interimFVs {
		var (
			interimElement = mod_strings.SplitN(b, dnSeparator, mod_slices.KVElements, mod_slices.FlagFilterEmpty|mod_slices.FlagTrimSpace)
		)
		switch {
		case len(interimElement) != mod_slices.KVElements:
			return mod_errors.EParse
		}

		for _, d := range interimElement {
			switch {
			case len(d) == 0:
				return mod_errors.EParse
			}
		}

		interim[a] = struct{ Field, Value string }{interimElement[0], interimElement[1]}
	}

	*r = interim

	return
}

func (r *attrDN) MarshalText() (outbound []byte, err error) {
	var (
		interim = make([]string, len(*r), len(*r))
	)
	for a, b := range *r {
		interim[a] = strings.Join([]string{b.Field, b.Value}, dnSeparator)
	}

	return []byte(strings.Join(interim, dnPathSeparator)), nil
}

func (r *attrDN) String() (outbound string) { return string(mod_errors.StripErr1(r.MarshalText())) }

func parseDN(inbound string) (outbound attrDN, err error) {
	var (
		interim = new(attrDN)
	)
	switch err = interim.UnmarshalText([]byte(inbound)); {
	case err != nil:
		return
	}

	return *interim, err
}

func (r *attrDN) parse(inbound string) (err error) {
	var (
		interim = new(attrDN)
	)

	switch err = interim.UnmarshalText([]byte(inbound)); {
	case err != nil:
		return
	}

	*r = *interim

	return
}
