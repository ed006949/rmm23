package mod_slices

// normalization flags.
const (
	FlagNone flag = 0
	FlagSort flag = 1 << iota
	FlagCompact
	FlagFilterEmpty
	FlagNormalize = ^flag(FlagNone)
)

const (
	KVElements = 2 // to honor the lint.mnd()
)
