package mod_bools

type errorNumber int

const (
	EINVAL errorNumber = iota
	ENODATA
)

var errorDescription = [...]string{
	EINVAL:  "invalid argument",
	ENODATA: "not data",
}

func (e errorNumber) Error() string        { return errorDescription[e] }
func (e errorNumber) Is(target error) bool { return e == target }
