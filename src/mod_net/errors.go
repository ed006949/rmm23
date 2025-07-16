package mod_net

type errorNumber int

const (
	ENODATA errorNumber = iota
)

var errorDescription = [...]string{
	ENODATA: "not data",
}

func (e errorNumber) Error() string        { return errorDescription[e] }
func (e errorNumber) Is(target error) bool { return e == target }
