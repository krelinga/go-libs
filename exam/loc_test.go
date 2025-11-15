package exam_test

import (
	"strings"
	"testing"

	"github.com/krelinga/go-libs/exam"
)

func TestLocString(t *testing.T) {
	tests := []struct {
		name     string
		loc      exam.Loc
		expected string
	}{
		{
			name:     "basic location",
			loc:      exam.Loc{File: "/path/to/file.go", Line: 42},
			expected: "/path/to/file.go:42",
		},
		{
			name:     "relative path",
			loc:      exam.Loc{File: "file.go", Line: 1},
			expected: "file.go:1",
		},
		{
			name:     "zero line",
			loc:      exam.Loc{File: "test.go", Line: 0},
			expected: "test.go:0",
		},
		{
			name:     "zero value",
			loc:      exam.Loc{},
			expected: "<uninitialized>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.loc.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestHere(t *testing.T) {
	// Call Here() and verify it captures the correct file and line
	loc1 := exam.Here()
	t.Log("Captured location:", loc1)

	// Check that the file contains "loc_test.go"
	if !strings.HasSuffix(loc1.File, "loc_test.go") {
		t.Errorf("Here() captured file %q, want file ending with 'loc_test.go'", loc1.File)
	}

	// The line should be positive
	if loc1.Line <= 0 {
		t.Errorf("Here() captured line %d, want positive line number", loc1.Line)
	}

	// Verify String() works with the captured location
	str := loc1.String()
	if !strings.Contains(str, "loc_test.go") {
		t.Errorf("String() = %q, want string containing 'loc_test.go'", str)
	}
	if !strings.Contains(str, ":") {
		t.Errorf("String() = %q, want string containing ':'", str)
	}

	// Test that consecutive calls to Here() return consecutive line numbers
	loc2 := exam.Here()
	loc3 := exam.Here()

	if loc2.File != loc3.File {
		t.Errorf("Here() calls returned different files: %q vs %q", loc2.File, loc3.File)
	}

	// loc3 should be exactly 1 line after loc2
	if loc3.Line != loc2.Line+1 {
		t.Errorf("Here() line numbers: loc2=%d, loc3=%d, want consecutive lines", loc2.Line, loc3.Line)
	}
}