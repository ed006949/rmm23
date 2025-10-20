package mod_db

import (
	"math"
	"strconv"
)

type attrEntryStatus int //

// entry status.
const (
	EntryStatusUnknown attrEntryStatus = iota
	EntryStatusLoaded
	EntryStatusCreated
	EntryStatusUpdated
	EntryStatusDeleted
	EntryStatusInvalid
	EntryStatusParsed
	EntryStatusSanitized
)
const (
	entryStatusReady = math.MaxInt
)

var (
	entryStatusNumber = map[attrEntryStatus]string{
		EntryStatusUnknown:   strconv.FormatInt(int64(EntryStatusUnknown), 10),
		EntryStatusLoaded:    strconv.FormatInt(int64(EntryStatusLoaded), 10),
		EntryStatusCreated:   strconv.FormatInt(int64(EntryStatusCreated), 10),
		EntryStatusUpdated:   strconv.FormatInt(int64(EntryStatusUpdated), 10),
		EntryStatusDeleted:   strconv.FormatInt(int64(EntryStatusDeleted), 10),
		EntryStatusInvalid:   strconv.FormatInt(int64(EntryStatusInvalid), 10),
		EntryStatusParsed:    strconv.FormatInt(int64(EntryStatusParsed), 10),
		EntryStatusSanitized: strconv.FormatInt(int64(EntryStatusSanitized), 10),
	}
	entryStatusString = map[attrEntryStatus]string{
		EntryStatusUnknown:   "unknown",
		EntryStatusLoaded:    "loaded",
		EntryStatusCreated:   "created",
		EntryStatusUpdated:   "updated",
		EntryStatusDeleted:   "deleted",
		EntryStatusInvalid:   "invalid",
		EntryStatusParsed:    "parsed",
		EntryStatusSanitized: "sanitized",
	}
)

func (r attrEntryStatus) Number() (outbound string) { return entryStatusNumber[r] }
func (r attrEntryStatus) String() (outbound string) { return entryStatusString[r] }
