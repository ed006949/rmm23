package mod_slices

// normalization flags.
const (
	FlagNone flagType = 0
	FlagSort flagType = 1 << iota
	FlagCompact
	FlagFilterEmpty
	FlagTrimSpace
	FlagNormalize = ^FlagNone
)

const (
	KVElements = 2 // to honor the lint.mnd()
)
