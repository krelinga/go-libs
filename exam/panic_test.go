package exam_test

import (
	"errors"
	"testing"

	"github.com/krelinga/go-libs/deep"
	"github.com/krelinga/go-libs/exam"
	"github.com/krelinga/go-libs/match"
)

func TestPanicWith(t *testing.T) {
	env := deep.NewEnv()
	err := errors.New("expected error panic")
	tests := []struct {
		name      string
		matcher   match.Matcher
		panicWith any
		wantOk    bool
		wantLog   match.Matcher
	}{
		{
			name:      "panic with matching string value",
			matcher:   match.Equal("expected panic"),
			panicWith: "expected panic",
			wantOk:    true,
			wantLog:   matchLogOk(),
		},
		{
			name:      "panic with non-matching string value",
			matcher:   match.Equal("expected panic"),
			panicWith: "unexpected panic",
			wantOk:    false,
			wantLog:   matchLogError(),
		},
		{
			name:      "no panic",
			matcher:   match.Equal("expected panic"),
			panicWith: nil,
			wantOk:    false,
			wantLog:   matchLogError(),
		},
		{
			name:      "panic with matching error value",
			matcher:   match.As[error](match.ErrorIs(err)),
			panicWith: err,
			wantOk:    true,
			wantLog:   matchLogOk(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				res := exam.PanicWith(e, env, tt.matcher, func() {
					if tt.panicWith != nil {
						panic(tt.panicWith)
					}
				})
				if res.Ok() != tt.wantOk {
					t.Errorf("PanicWith().Ok() = %v, want %v", res.Ok(), tt.wantOk)
				}
			})
			result := tt.wantLog.Match(env, match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("PanicWith() log = %s", deep.Format(env, log))
			}
		})
	}
}
