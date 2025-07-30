package mod_db

type qs struct {
	_Q  []_Q
	_OF []entryFieldName
}

type _Q struct {
	_F entryFieldName
	_V string
}
