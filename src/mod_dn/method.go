package mod_dn

import (
	"strings"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
)

func (r *DN) UnmarshalText(inbound []byte) (err error) {
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

		interim[a] = mod_strings.KV{interimElement[0], interimElement[1]}
	}

	r.dn = interim

	return
}

func (r *DN) MarshalText() (outbound []byte, err error) {
	var (
		interim = make([]string, len(r.dn), len(r.dn))
	)
	for a, b := range r.dn {
		switch {
		case len(b.Key) == 0 || len(b.Value) == 0:
			return nil, mod_errors.EParse
		}

		interim[a] = strings.Join([]string{b.Key, b.Value}, mod_strings.DNKVSeparator)
	}

	return []byte(strings.Join(interim, mod_strings.DNPathSeparator)), nil
}

func (r *DN) String() (outbound string) { return string(mod_errors.StripErr1(r.MarshalText())) }

func (r *DN) Parse(inbound string) (err error) {
	var (
		interim = new(DN)
	)

	switch err = interim.UnmarshalText([]byte(inbound)); {
	case err != nil:
		return
	}

	*r = *interim

	return
}
