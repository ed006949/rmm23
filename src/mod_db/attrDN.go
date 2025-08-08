package mod_db

import (
	"strings"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
)

type attrDN struct {
	dn []mod_strings.KV
}

func (r *attrDN) UnmarshalText(inbound []byte) (err error) {
	var (
		interimFVs = mod_strings.Split(string(inbound), mod_strings.DNPathSeparator, mod_slices.FlagFilterEmpty|mod_slices.FlagTrimSpace)
		interim    = make([]mod_strings.KV, len(interimFVs))
	)
	for a, b := range interimFVs {
		var (
			interimElement = mod_strings.SplitN(b, mod_strings.DNKVSeparator, mod_slices.KVElements, mod_slices.FlagFilterEmpty|mod_slices.FlagTrimSpace)
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

		interim[a] = mod_strings.KV{interimElement[0], interimElement[1]}
	}

	r.dn = interim

	return
}

func (r *attrDN) MarshalText() (outbound []byte, err error) {
	var (
		interim = make([]string, len(r.dn), len(r.dn))
	)
	for a, b := range r.dn {
		interim[a] = strings.Join([]string{b.Key, b.Value}, mod_strings.DNKVSeparator)
	}

	return []byte(strings.Join(interim, mod_strings.DNPathSeparator)), nil
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
