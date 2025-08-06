package mod_db

import (
	"encoding/json"
	"strings"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
)

type attrDN []struct{ Field, Value string }
type attrDNs []*attrDN

func (r *attrDN) MarshalJSON() (outbound []byte, err error) { return json.Marshal(r.String()) }

func (r *attrDN) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim string
	)
	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	switch err = r.parse(interim); {
	case err != nil:
		return
	}

	return
}

func (r *attrDN) UnmarshalLDAPAttr(values []string) (err error) {
	var (
		interim attrDN
	)
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		switch err = interim.parse(value); {
		case err != nil:
			return
		}

		*r = interim

		return // return only first value
	}

	return
}

func (r *attrDNs) UnmarshalLDAPAttr(values []string) (err error) {
	var (
		interim attrDNs
	)
	switch err = interim.parse(mod_slices.StringsNormalize(values, mod_slices.FlagNormalize)); {
	case err != nil:
		return
	}

	*r = interim

	return
}

func (r *attrDN) String() (outbound string) {
	var (
		interim = make([]string, len(*r), len(*r))
	)
	for a, b := range *r {
		interim[a] = strings.Join([]string{b.Field, b.Value}, "=")
	}

	return strings.Join(interim, ",")
}

func (r *attrDNs) String() (outbound []string) {
	for _, b := range *r {
		outbound = append(outbound, b.String())
	}

	return
}

func (r *attrDN) parse(inbound string) (err error) {
	var (
		interim attrDN
	)
	switch interim, err = parseDN(inbound); {
	case err != nil:
		return
	}

	*r = interim

	return
}

func parseDN(inbound string) (outbound attrDN, err error) {
	var (
		interimFVs = mod_slices.SplitString(inbound, ",", mod_slices.FlagFilterEmpty|mod_slices.FlagTrimSpace)
		interim    = make(attrDN, len(interimFVs), len(interimFVs))
	)
	for a, b := range interimFVs {
		var (
			interimElement = mod_slices.SplitStringN(b, "=", mod_slices.KVElements, mod_slices.FlagFilterEmpty|mod_slices.FlagTrimSpace)
		)
		switch {
		case len(interimElement) != mod_slices.KVElements:
			return nil, mod_errors.EParse
		}

		for _, d := range interimElement {
			switch {
			case len(d) == 0:
				return nil, mod_errors.EParse
			}
		}

		interim[a] = struct{ Field, Value string }{interimElement[0], interimElement[1]}
	}

	return interim, nil
}

func (r *attrDNs) parse(inbound []string) (err error) {
	var (
		interim = make(attrDNs, len(inbound), len(inbound))
	)
	for a, b := range inbound {
		interim[a] = new(attrDN)
		switch err = interim[a].parse(b); {
		case err != nil:
			return mod_errors.EParse
		}
	}

	*r = interim

	return
}
