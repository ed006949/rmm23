package mod_slices

// normalization flags
const (
	FlagSort flag = 1 << iota
	FlagCompact
	FlagFilterEmpty
	FlagNormalize = ^flag(0)
)
