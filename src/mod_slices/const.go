package mod_slices

// normalization flags
const (
	FlagNone flag = 0
	FlagSort flag = 1 << iota
	FlagCompact
	FlagFilterEmpty
	FlagNormalize = ^flag(FlagNone)
)
