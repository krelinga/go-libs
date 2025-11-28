package match_test

import (
	"errors"
	"testing"

	"github.com/krelinga/go-libs/deep"
	"github.com/krelinga/go-libs/match"
)

func TestErrorIs(t *testing.T) {
	err := errors.New("file not found")
	tests := []struct {
		name    string
		in      match.Vals
		target  error
		wantMatch bool
	} {
		{
			name:     "error matches target",
			in:       match.NewVals1(err),
			target:   err,
			wantMatch: true,
		},
		{
			name:     "error does not match target",
			in:       match.NewVals1(errors.New("different error")),
			target:   err,
			wantMatch: false,
		},
		// {
		// 	name:     "nil error does not match target",
		// 	in:       match.NewVals1[error](nil),
		// 	target:   err,
		// 	wantMatch: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := match.ErrorIs(tt.target)
			result := matcher.Match(deep.NewEnv(), tt.in)
			if result.Matched() != tt.wantMatch {
				t.Errorf("ErrorIs().Match() = %v, want %v", result.Matched(), tt.wantMatch)
			}
		})
	}
}