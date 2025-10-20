package mod_db

import (
	"math"
	"strconv"

	"rmm23/src/mod_strings"
)

type attrEntryStatus int //

// entry status.
const (
	EntryStatusUnknown attrEntryStatus = iota
	EntryStatusLoad
	EntryStatusCreate
	EntryStatusUpdate
	EntryStatusDelete
	EntryStatusInvalid
	EntryStatusParse
	EntryStatusSanitize
)
const (
	EntryStatusReady = math.MaxInt
)

var (
	entryStatusMap = []mod_strings.MDMap{
		EntryStatusUnknown: {
			Number: strconv.FormatInt(int64(EntryStatusUnknown), 10),
			String: "unknown",
		},
		EntryStatusLoad: {
			Number: strconv.FormatInt(int64(EntryStatusLoad), 10),
			String: "load",
		},
		EntryStatusCreate: {
			Number: strconv.FormatInt(int64(EntryStatusCreate), 10),
			String: "create",
		},
		EntryStatusUpdate: {
			Number: strconv.FormatInt(int64(EntryStatusUpdate), 10),
			String: "update",
		},
		EntryStatusDelete: {
			Number: strconv.FormatInt(int64(EntryStatusDelete), 10),
			String: "delete",
		},
		EntryStatusInvalid: {
			Number: strconv.FormatInt(int64(EntryStatusInvalid), 10),
			String: "invalid",
		},
		EntryStatusParse: {
			Number: strconv.FormatInt(int64(EntryStatusParse), 10),
			String: "parse",
		},
		EntryStatusSanitize: {
			Number: strconv.FormatInt(int64(EntryStatusSanitize), 10),
			String: "sanitize",
		},
	}
)

func (r attrEntryStatus) Number() (outbound string) { return entryStatusMap[r].Number }
func (r attrEntryStatus) String() (outbound string) { return entryStatusMap[r].String }
