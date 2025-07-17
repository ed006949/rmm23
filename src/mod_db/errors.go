package mod_db

const (
	EDocExist errorNumber = iota
)

var errorDescription = [...]string{
	EDocExist: "Document already exists",
}

type errorNumber int

func (e errorNumber) Error() string        { return errorDescription[e] }
func (e errorNumber) Is(target error) bool { return e == target }
