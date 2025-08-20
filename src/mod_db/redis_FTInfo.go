package mod_db

import (
	"encoding/json/v2"

	"rmm23/src/mod_reflect"
	"rmm23/src/mod_strings"
)

// ftInfo mirrors the logical FT.INFO reply (RediSearch ≥ 2.0).
type ftInfo struct {
	IndexName       string                `json:"index_name"`
	IndexOptions    map[string]any        `json:"index_options"`
	IndexDefinition ftInfoIndexDefinition `json:"index_definition"`
	Attributes      ftInfoAttributes      `json:"attributes"`

	/* Counters */
	NumDocs    int64 `json:"num_docs"`
	MaxDocID   int64 `json:"max_doc_id"`
	NumTerms   int64 `json:"num_terms"`
	NumRecords int64 `json:"num_records"`

	/* Size / memory */
	InvertedSizeMB           float64 `json:"inverted_sz_mb"`
	TotalInvertedIndexBlocks int64   `json:"total_inverted_index_blocks"`
	OffsetVectorsSizeMB      float64 `json:"offset_vectors_sz_mb"`
	DocTableSizeMB           float64 `json:"doc_table_size_mb"`
	SortableValuesSizeMB     float64 `json:"sortable_values_size_mb"`
	KeyTableSizeMB           float64 `json:"key_table_size_mb"`
	RecordsPerDocAvg         float64 `json:"records_per_doc_avg"`
	BytesPerRecordAvg        float64 `json:"bytes_per_record_avg"`
	OffsetsPerTermAvg        float64 `json:"offsets_per_term_avg"`
	OffsetBitsPerRecordAvg   string  `json:"offset_bits_per_record_avg"` // can be "NaN"

	/* Indexing progress */
	HashIndexingFailures int64   `json:"hash_indexing_failures"`
	Indexing             int64   `json:"indexing"` // 0 | 1
	PercentIndexed       float64 `json:"percent_indexed"`

	/* Sub-objects */
	GCStats     ftInfoGCStats     `json:"gc_stats"`
	CursorStats ftInfoCursorStats `json:"cursor_stats"`
}

type ftInfoIndexDefinition struct {
	KeyType      string   `json:"key_type"`      // "JSON"
	Prefixes     []string `json:"prefixes"`      // ["certificate:"]
	DefaultScore float64  `json:"default_score"` // 1
}

type ftInfoAttribute struct {
	Identifier string `json:"identifier"`          // $.status, $.uuid, …
	Attribute  string `json:"attribute"`           // alias (status, uuid …)
	Type       string `json:"type"`                // NUMERIC | TAG | …
	Separator  string `json:"separator,omitempty"` // TAG only
}

type ftInfoGCStats struct {
	BytesCollected       int64  `json:"bytes_collected"`
	TotalMsRun           int64  `json:"total_ms_run"`
	TotalCycles          int64  `json:"total_cycles"`
	AverageCycleTimeMs   string `json:"average_cycle_time_ms"` // may be "NaN"
	LastRunTimeMs        int64  `json:"last_run_time_ms"`
	GCNumericTreesMissed int64  `json:"gc_numeric_trees_missed"`
	GCBlocksDenied       int64  `json:"gc_blocks_denied"`
}

type ftInfoCursorStats struct {
	GlobalIdle    int64 `json:"global_idle"`
	GlobalTotal   int64 `json:"global_total"`
	IndexCapacity int64 `json:"index_capacity"`
	IndexTotal    int64 `json:"index_total"`
}

type ftInfoAttributes map[mod_strings.EntryFieldName]*ftInfoAttribute

func (r *ftInfoAttributes) UnmarshalJSON(data []byte) (err error) {
	mod_reflect.MakeMapIfNil(r)

	var (
		interim []*ftInfoAttribute
	)
	switch err = json.Unmarshal(data, &interim); {
	case err != nil:
		return
	}

	for _, b := range interim {
		(*r)[mod_strings.EntryFieldName(b.Attribute)] = b
	}

	return
}
