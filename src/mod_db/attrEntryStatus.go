package mod_db

import (
	"math"
	"strconv"

	"rmm23/src/mod_strings"
)

type attrEntryStatus int //

// entry status.
const (
	entryStatusUnknown attrEntryStatus = iota
	entryStatusLoad
	entryStatusCreate
	entryStatusUpdate
	entryStatusDelete
	entryStatusInvalid
	entryStatusParse
	entryStatusSanitize
)
const (
	entryStatusReady = math.MaxInt
)

var (
	entryStatusMap = []mod_strings.MDMap{
		entryStatusUnknown: {
			Number: strconv.FormatInt(int64(entryStatusUnknown), 10),
			String: "unknown",
		},
		entryStatusLoad: {
			Number: strconv.FormatInt(int64(entryStatusLoad), 10),
			String: "load",
		},
		entryStatusCreate: {
			Number: strconv.FormatInt(int64(entryStatusCreate), 10),
			String: "create",
		},
		entryStatusUpdate: {
			Number: strconv.FormatInt(int64(entryStatusUpdate), 10),
			String: "update",
		},
		entryStatusDelete: {
			Number: strconv.FormatInt(int64(entryStatusDelete), 10),
			String: "delete",
		},
		entryStatusInvalid: {
			Number: strconv.FormatInt(int64(entryStatusInvalid), 10),
			String: "invalid",
		},
		entryStatusParse: {
			Number: strconv.FormatInt(int64(entryStatusParse), 10),
			String: "parse",
		},
		entryStatusSanitize: {
			Number: strconv.FormatInt(int64(entryStatusSanitize), 10),
			String: "sanitize",
		},
	}
)

func (r attrEntryStatus) Number() (outbound string) { return entryStatusMap[r].Number }
func (r attrEntryStatus) String() (outbound string) { return entryStatusMap[r].String }
