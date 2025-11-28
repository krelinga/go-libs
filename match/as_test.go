package match_test

import (
	"errors"
	"testing"

	"github.com/krelinga/go-libs/deep"
	"github.com/krelinga/go-libs/match"
)

func TestAs(t *testing.T) {
	type MyInt int
	tests := []struct {
		name    string
		vals     match.Vals
		matcher match.Matcher
		wantMatch bool
	}{
		{
			name:     "int type conversion",
			vals:     match.NewVals1[int32](42),
			matcher:  match.As[int64](match.Equal[int64](42)),
			wantMatch: true,
		},
		{
			name: "derived type conversion",
			vals: match.NewVals1(7),
			matcher: match.As[MyInt](match.Equal(MyInt(7))),
			wantMatch: true,
		},
		{
			name: "interface type conversion",
			vals: match.NewVals1[any](errors.New("test error")),
			matcher: match.As[error](match.Equal(errors.New("test error"))),
			wantMatch: true,
		},
		{
			name: "failed type conversion",
			vals: match.NewVals1("not an int"),
			matcher: match.As[int](match.Equal(0)),
			wantMatch: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.matcher.Match(deep.NewEnv(), tt.vals)
			if result.Matched() != tt.wantMatch {
				t.Errorf("As().Match() = %v, want %v", result.Matched(), tt.wantMatch)
			}
		})
	}
}
