package mod_strings

import (
	"reflect"
	"testing"

	"rmm23/src/mod_slices"
)

func TestToStrings(t *testing.T) {
	tests := []struct {
		name     string
		inbound  []any
		flag     mod_slices.FlagType
		expected []string
	}{
		{
			name:     "simple strings",
			inbound:  []any{"a", "b", "c"},
			flag:     0,
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "mixed types",
			inbound:  []any{"a", 1, true},
			flag:     0,
			expected: []string{"a", "1", "true"},
		},
		{
			name:     "with trim space flag",
			inbound:  []any{" a ", " b ", " c "},
			flag:     mod_slices.FlagTrimSpace,
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with filter empty flag",
			inbound:  []any{"a", "", "c"},
			flag:     mod_slices.FlagFilterEmpty,
			expected: []string{"a", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToStrings(tt.inbound, tt.flag)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ToStrings() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestJoinAny(t *testing.T) {
	tests := []struct {
		name     string
		inbound  []any
		sep      string
		flag     mod_slices.FlagType
		expected string
	}{
		{
			name:     "simple join",
			inbound:  []any{"a", "b", "c"},
			sep:      ",",
			flag:     0,
			expected: "a,b,c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JoinAny(tt.inbound, tt.sep, tt.flag)
			if got != tt.expected {
				t.Errorf("JoinAny() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		inbound  string
		sep      string
		flag     mod_slices.FlagType
		expected []string
	}{
		{
			name:     "simple split",
			inbound:  "a,b,c",
			sep:      ",",
			flag:     0,
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "split with trim and filter",
			inbound:  " a , , c ",
			sep:      ",",
			flag:     mod_slices.FlagTrimSpace | mod_slices.FlagFilterEmpty,
			expected: []string{"a", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split(tt.inbound, tt.sep, tt.flag)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Split() = %v, want %v", got, tt.expected)
			}
		})
	}
}
