package mod_db

import (
	"math"
	"strconv"
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
	entryStatusNumber = map[attrEntryStatus]string{
		EntryStatusUnknown:  strconv.FormatInt(int64(EntryStatusUnknown), 10),
		EntryStatusLoad:     strconv.FormatInt(int64(EntryStatusLoad), 10),
		EntryStatusCreate:   strconv.FormatInt(int64(EntryStatusCreate), 10),
		EntryStatusUpdate:   strconv.FormatInt(int64(EntryStatusUpdate), 10),
		EntryStatusDelete:   strconv.FormatInt(int64(EntryStatusDelete), 10),
		EntryStatusInvalid:  strconv.FormatInt(int64(EntryStatusInvalid), 10),
		EntryStatusParse:    strconv.FormatInt(int64(EntryStatusParse), 10),
		EntryStatusSanitize: strconv.FormatInt(int64(EntryStatusSanitize), 10),
	}
	entryStatusString = map[attrEntryStatus]string{
		EntryStatusUnknown:  "unknown",
		EntryStatusLoad:     "load",
		EntryStatusCreate:   "create",
		EntryStatusUpdate:   "update",
		EntryStatusDelete:   "delete",
		EntryStatusInvalid:  "invalid",
		EntryStatusParse:    "parse",
		EntryStatusSanitize: "sanitize",
	}
)

func (r attrEntryStatus) Number() (outbound string) { return entryStatusNumber[r] }
func (r attrEntryStatus) String() (outbound string) { return entryStatusString[r] }
