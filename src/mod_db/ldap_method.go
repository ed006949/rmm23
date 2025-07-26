package mod_db

import (
	"github.com/google/uuid"
)

func (r *AttrDN) String() string { return string(*r) }

func (r *AttrUUID) String() string { return uuid.UUID(*r).String() }
func (r *AttrUUID) Entry() string  { return entryDocIDHeader + r.String() }
