package mod_ldap

const (
	EAnonymousBind errorNumber = iota
	ENoConn
	ENotStruct
	ENotPtr
	EUnknownType
	EParse
)

var errorDescription = [...]string{
	EAnonymousBind: "anonymous bind",
	ENoConn:        "no connection",
	ENotStruct:     "not a struct",
	ENotPtr:        "not a pointer",
	EUnknownType:   "unknown type",
	EParse:         "parse error",
}

type errorNumber int

func (e errorNumber) Error() string        { return errorDescription[e] }
func (e errorNumber) Is(target error) bool { return e == target }
