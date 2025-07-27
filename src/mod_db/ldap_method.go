package mod_db

import (
	"github.com/google/uuid"
)

func (r *attrDN) String() string { return string(*r) }

func (r *attrUUID) String() string { return uuid.UUID(*r).String() }
func (r *attrUUID) Entry() string  { return entryDocIDHeader + r.String() }
