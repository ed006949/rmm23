package mod_strings

type EntryFieldName string
type EntryFieldMap map[EntryFieldName]string

type FVs []FV
type FV struct {
	Field EntryFieldName
	Value string
}

type KVs []KV
type KV struct {
	Key   string
	Value string
}
