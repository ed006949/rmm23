package mod_db

func (r entryFieldName) String() (outbound string)         { return string(r) }
func (r entryFieldName) FieldName() (outbound string)      { return "$." + string(r) }
func (r entryFieldName) FieldNameSlice() (outbound string) { return "$." + string(r) + "[*]" }
