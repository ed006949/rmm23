package mod_strings

func (r EntryFieldName) String() (outbound string)         { return string(r) }
func (r EntryFieldName) FieldName() (outbound string)      { return JSONPathHeader + r.String() }
func (r EntryFieldName) FieldNameSlice() (outbound string) { return r.FieldName() + "[*]" }
