package mod_db

import (
	"encoding/json"
	"fmt"
	"strings"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
)

const (
	dnSeparator     = "="
	dnPathSeparator = ","
)

type attrDN []struct{ Field, Value string }
type attrDNs []*attrDN

func (r *attrDN) UnmarshalText(text []byte) error {
	var (
		s = string(text)
	)
	switch {
	case len(s) == 0:
		*r = nil

		return nil
	}

	var (
		// Split the DN on commas, respecting simple RFC4514 escaping
		dnParts = make([]struct{ Field, Value string }, len(text), len(text))
	)

	fields := strings.Split(s, ",")
	for _, f := range fields {
		kv := strings.SplitN(strings.TrimSpace(f), mod_slices.KVSeparator, mod_slices.KVElements)
		if len(kv) != mod_slices.KVElements {
			return fmt.Errorf("attrDN: invalid DN component: %q", f)
		}

		dnParts = append(dnParts, struct {
			Field, Value string
		}{
			Field: kv[0],
			Value: kv[1],
		})
	}

	*r = dnParts

	return nil
}
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

func (r *attrDN) String() (outbound string) {
	var (
		interim = make([]string, len(*r), len(*r))
	)
	for a, b := range *r {
		interim[a] = strings.Join([]string{b.Field, b.Value}, dnSeparator)
	}

	return strings.Join(interim, dnPathSeparator)
}

func (r *attrDNs) String() (outbound []string) {
	for _, b := range *r {
		outbound = append(outbound, b.String())
	}

	return
}

func parseDN(inbound string) (outbound *attrDN, err error) {
	var (
		interim = new(attrDN)
	)
	switch err = interim.parse(inbound); {
	case err != nil:
		return
	}

	return interim, err
}

func (r *attrDN) parse(inbound string) (err error) {
	var (
		interimFVs = mod_strings.Split(inbound, dnPathSeparator, mod_slices.FlagFilterEmpty|mod_slices.FlagTrimSpace)
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
