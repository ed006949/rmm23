package mod_errors

const (
	EUnwilling errorNumber = iota
	EORPHANED
	EDUPDATA
	EEXIST
	ENOTFOUND
	EINVAL
	ENODATA
	ENEDATA
	ENOCONF
	EUEDATA
	EINVALRESPONSE
	EAnonymousBind
	ENoConn

	ENotStruct
	ENotPtr
	ENotStructOrPtrStruct
	ENotPtrOrNil

	EUnknownType
	EParse

	EX509ParsePrivKey
	EPEMNoDataKey
	EPEMNoDataCert
	EUnknownPrivKeyType
	EUnknownPubKeyAlgo
	EPrivKeyType
	EPrivKeySize
	ETypeMismatchPrivKeyPubKey
	EMismatchPrivKeyPubKey

	ECom
	EComSet
	EComSetDomAdm
	EComSetDomSetAdm

	ETagMultiType
	ETagUnknown

	EUnknownScheme
	EDocExist
	EIndexExist
)

var errorDescription = [...]string{
	EUnwilling:                 "unwilling to perform",
	EORPHANED:                  "orphaned entry",
	EDUPDATA:                   "duplicate data",
	EEXIST:                     "already exists",
	ENOTFOUND:                  "not found",
	EINVAL:                     "invalid argument",
	ENODATA:                    "not data",
	ENEDATA:                    "not enough data",
	ENOCONF:                    "no config",
	EUEDATA:                    "unexpected data",
	EINVALRESPONSE:             "invalid response",
	EAnonymousBind:             "anonymous bind",
	ENoConn:                    "no connection",
	ENotStruct:                 "not a struct",
	ENotPtr:                    "not a pointer",
	ENotStructOrPtrStruct:      "not a struct or a pointer to a struct",
	ENotPtrOrNil:               "not a pointer or nil",
	EUnknownType:               "unknown type",
	EParse:                     "parse error",
	EX509ParsePrivKey:          "x509: failed to parse private key",
	EPEMNoDataKey:              "PEM: failed to find any PRIVATE KEY data",
	EPEMNoDataCert:             "PEM: failed to find any CERTIFICATE data",
	EUnknownPrivKeyType:        "unknown private key type",
	EUnknownPubKeyAlgo:         "unknown public key algorithm",
	EPrivKeyType:               "wrong private key type",
	EPrivKeySize:               "wrong private key size",
	ETypeMismatchPrivKeyPubKey: "private key type does not match public key type",
	EMismatchPrivKeyPubKey:     "private key does not match public key",
	ECom:                       "unknown command",
	EComSet:                    "unknown command set",
	EComSetDomAdm:              "unknown Domain Administration command",
	EComSetDomSetAdm:           "unknown Domain Set Administration command",
	ETagMultiType:              "multiple types",
	ETagUnknown:                "unknown tag",
	EUnknownScheme:             "unknown scheme",
	EDocExist:                  "Document already exists",
	EIndexExist:                "Index already exists",
}
