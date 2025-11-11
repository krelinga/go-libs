package exam_test

import (
	"testing"

	"github.com/krelinga/go-deep"
	"github.com/krelinga/go-deep/exam"
	"github.com/krelinga/go-deep/match"
)

func TestNil(t *testing.T) {
	tests := []struct {
		name    string
		a       any
		want    bool
		matcher match.Matcher
	}{
		{
			name: "a is nil",
			a:    nil,
			want: true,
			matcher: matchLogOk(),
		},
		{
			name: "a is not nil",
			a:    5,
			want: false,
			matcher: matchLogError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				if ok := exam.Nil(e, deep.NewEnv(), tt.a).Ok(); ok != tt.want {
					t.Errorf("Nil(%v) = %v, want %v", tt.a, ok, tt.want)
				}
			})
			result := tt.matcher.Match(deep.NewEnv(), match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("Nil() log = %s", deep.Format(deep.NewEnv(), log))
			}
		})
	}
}

func TestNotNil(t *testing.T) {
	tests := []struct {
		name    string
		a       any
		want    bool
		matcher match.Matcher
	}{
		{
			name: "a is not nil",
			a:    5,
			want: true,
			matcher: matchLogOk(),
		},
		{
			name: "a is nil",
			a:    nil,
			want: false,
			matcher: matchLogError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				if ok := exam.NotNil(e, deep.NewEnv(), tt.a).Ok(); ok != tt.want {
					t.Errorf("NotNil(%v) = %v, want %v", tt.a, ok, tt.want)
				}
			})
			result := tt.matcher.Match(deep.NewEnv(), match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("NotNil() log = %s", deep.Format(deep.NewEnv(), log))
			}
		})
	}
}
