package mod_slices

// Has checks if a specific flag is set within the current flags bitmask.
func (r *FlagType) Has(flag FlagType) bool { return (*r & flag) != 0 }
