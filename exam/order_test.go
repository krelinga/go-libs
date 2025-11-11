package exam_test

import (
	"testing"

	"github.com/krelinga/go-deep"
	"github.com/krelinga/go-deep/exam"
	"github.com/krelinga/go-deep/match"
)

func TestLessThan(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    bool
		matcher match.Matcher
	}{
		{
			name: "a less than b",
			a:    3,
			b:    5,
			want: true,
			matcher: matchLogOk(),
		},
		{
			name: "a equal to b",
			a:    5,
			b:    5,
			want: false,
			matcher: matchLogError(),
		},
		{
			name: "a greater than b",
			a:    7,
			b:    5,
			want: false,
			matcher: matchLogError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				env := deep.NewEnv()
				if ok := exam.LessThan(e, env, tt.a, tt.b).Ok(); ok != tt.want {
					t.Errorf("LessThan(%d, %d) = %v, want %v", tt.a, tt.b, ok, tt.want)
				}
			})
			result := tt.matcher.Match(deep.NewEnv(), match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("LessThan() log = %s", deep.Format(deep.NewEnv(), log))
			}
		})
	}
}

func TestGreaterThan(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    bool
		matcher match.Matcher
	}{
		{
			name: "a greater than b",
			a:    7,
			b:    5,
			want: true,
			matcher: matchLogOk(),
		},
		{
			name: "a equal to b",
			a:    5,
			b:    5,
			want: false,
			matcher: matchLogError(),
		},
		{
			name: "a less than b",
			a:    3,
			b:    5,
			want: false,
			matcher: matchLogError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				env := deep.NewEnv()
				if ok := exam.GreaterThan(e, env, tt.a, tt.b).Ok(); ok != tt.want {
					t.Errorf("GreaterThan(%d, %d) = %v, want %v", tt.a, tt.b, ok, tt.want)
				}
			})
			result := tt.matcher.Match(deep.NewEnv(), match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("LessThan() log = %s", deep.Format(deep.NewEnv(), log))
			}
		})
	}
}

func TestLessThanOrEqual(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    bool
		matcher match.Matcher
	}{
		{
			name: "a less than b",
			a:    3,
			b:    5,
			want: true,
			matcher: matchLogOk(),
		},
		{
			name: "a equal to b",
			a:    5,
			b:    5,
			want: true,
			matcher: matchLogOk(),
		},
		{
			name: "a greater than b",
			a:    7,
			b:    5,
			want: false,
			matcher: matchLogError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				env := deep.NewEnv()
				if ok := exam.LessThanOrEqual(e, env, tt.a, tt.b).Ok(); ok != tt.want {
					t.Errorf("LessThanOrEqual(%d, %d) = %v, want %v", tt.a, tt.b, ok, tt.want)
				}
			})
			result := tt.matcher.Match(deep.NewEnv(), match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("LessThanOrEqual() log = %s", deep.Format(deep.NewEnv(), log))
			}
		})
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    bool
		matcher match.Matcher
	}{
		{
			name: "a greater than b",
			a:    7,
			b:    5,
			want: true,
			matcher: matchLogOk(),
		},
		{
			name: "a equal to b",
			a:    5,
			b:    5,
			want: true,
			matcher: matchLogOk(),
		},
		{
			name: "a less than b",
			a:    3,
			b:    5,
			want: false,
			matcher: matchLogError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				env := deep.NewEnv()
				if ok := exam.GreaterThanOrEqual(e, env, tt.a, tt.b).Ok(); ok != tt.want {
					t.Errorf("GreaterThanOrEqual(%d, %d) = %v, want %v", tt.a, tt.b, ok, tt.want)
				}
			})
			result := tt.matcher.Match(deep.NewEnv(), match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("GreaterThanOrEqual() log = %s", deep.Format(deep.NewEnv(), log))
			}
		})
	}
}
