package mod_slices

// has checks if a specific flag is set within the current flags bitmask.
func (r *flagType) has(flag flagType) bool { return (*r & flag) != 0 }
