package mod_strings

type EntryFieldName string

type KV struct {
	Key   string
	Value string
}

type FV struct {
	Field EntryFieldName
	Value string
}
