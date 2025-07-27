package mod_db

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (r *attrDN) String() (outbound string) { return string(*r) }
func (r *attrTimestamp) String() (outbound string) {
	return strconv.FormatInt(time.Time(*r).Unix(), 10)
}
func (r *attrUUID) String() (outbound string) { return uuid.UUID(*r).String() }
func (r *attrUUID) Entry() (outbound string)  { return entryDocIDHeader + r.String() }
