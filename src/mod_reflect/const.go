package mod_reflect

const (
	UnknownPointer = iota
	PointerToScalar
	PointerToSlice
	SliceOfPointers
	SliceOfValues
)
