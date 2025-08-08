package mod_db

import (
	"rmm23/src/mod_bytes"
	"rmm23/src/mod_reflect"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
)

type attrLabeledURIs map[string]string //

func (r *attrLabeledURIs) UnmarshalText(inbound []byte) (err error) {
	switch value := mod_bytes.SplitN(inbound, []byte(mod_strings.LURISeparator), mod_slices.KVElements, mod_slices.FlagNormalize); len(value) {
	case 1:
		mod_reflect.MakeMapIfNil(r)
		(*r)[string(value[0])] = ""
	case mod_slices.KVElements:
		mod_reflect.MakeMapIfNil(r)
		(*r)[string(value[0])] = string(value[1])
	}

	return
}
