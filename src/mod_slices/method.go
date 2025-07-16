package mod_slices

// has checks if a specific flag is set within the current flags bitmask.
func (r *flag) has(flag flag) bool { return (*r & flag) != 0 }
