package mod_slices

// normalization flags.
const (
	FlagNone FlagType = 0
	FlagSort FlagType = 1 << iota
	FlagCompact
	FlagFilterEmpty
	FlagTrimSpace
	FlagNormalize = ^FlagNone
)

const (
	KVElements  = 2 // to honor the lint.mnd()
	KVSeparator = "="
)
