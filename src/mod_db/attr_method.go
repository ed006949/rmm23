package mod_db

import (
	"github.com/google/uuid"
)

func (r *attrDN) String() (outbound string) { return string(*r) }

func (r *attrUUID) String() (outbound string) { return uuid.UUID(*r).String() }

func (r *attrUUID) Entry() (outbound string) { return entryDocIDHeader + r.String() }
