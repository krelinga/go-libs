package exam_test

import (
	"testing"

	"github.com/krelinga/go-deep/deep"
	"github.com/krelinga/go-deep/exam"
	"github.com/krelinga/go-deep/match"
)

func TestEqual(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    bool
		matcher match.Matcher
	}{
		{
			name:    "a equal to b",
			a:       5,
			b:       5,
			want:    true,
			matcher: matchLogOk(),
		},
		{
			name:    "a not equal to b",
			a:       5,
			b:       6,
			want:    false,
			matcher: matchLogError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				if ok := exam.Equal(e, deep.NewEnv(), tt.a, tt.b).Ok(); ok != tt.want {
					t.Errorf("Equal(%v, %v) = %v, want %v", tt.a, tt.b, ok, tt.want)
				}
			})
			result := tt.matcher.Match(deep.NewEnv(), match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("Equal() log = %s", deep.Format(deep.NewEnv(), log))
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    bool
		matcher match.Matcher
	}{
		{
			name:    "a not equal to b",
			a:       5,
			b:       6,
			want:    true,
			matcher: matchLogOk(),
		},
		{
			name:    "a equal to b",
			a:       5,
			b:       5,
			want:    false,
			matcher: matchLogError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				if ok := exam.NotEqual(e, deep.NewEnv(), tt.a, tt.b).Ok(); ok != tt.want {
					t.Errorf("NotEqual(%v, %v) = %v, want %v", tt.a, tt.b, ok, tt.want)
				}
			})
			result := tt.matcher.Match(deep.NewEnv(), match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("NotEqual() log = %s", deep.Format(deep.NewEnv(), log))
			}
		})
	}
}
