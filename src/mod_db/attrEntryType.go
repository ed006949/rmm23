package mod_db

import (
	"strconv"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_strings"
)

type attrEntryType int //

const (
	entryTypeEmpty attrEntryType = iota
	entryTypeDomain
	entryTypeGroup
	entryTypeUser
	entryTypeHost
)

var (
	entryTypeMap = []mod_strings.MDMap{
		entryTypeEmpty: {
			Number: strconv.FormatInt(int64(entryTypeEmpty), 10),
			String: "",
		},
		entryTypeDomain: {
			Number: strconv.FormatInt(int64(entryTypeDomain), 10),
			String: "domain",
		},
		entryTypeGroup: {
			Number: strconv.FormatInt(int64(entryTypeGroup), 10),
			String: "group",
		},
		entryTypeUser: {
			Number: strconv.FormatInt(int64(entryTypeUser), 10),
			String: "user",
		},
		entryTypeHost: {
			Number: strconv.FormatInt(int64(entryTypeHost), 10),
			String: "host",
		},
	}

	entryTypeID = func() (outbound map[string]attrEntryType) {
		outbound = make(map[string]attrEntryType, len(entryTypeMap))
		for a, b := range entryTypeMap {
			outbound[b.String] = attrEntryType(a)
		}

		return
	}()
)

func (r attrEntryType) Number() (outbound string) { return entryTypeMap[r].Number }
func (r attrEntryType) String() (outbound string) { return entryTypeMap[r].String }

func (r *attrEntryType) Parse(inbound string) (err error) {
	switch value, ok := entryTypeID[inbound]; {
	case !ok:
		return mod_errors.EUnknownType
	default:
		*r = value

		return
	}
}
