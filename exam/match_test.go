package exam_test

import (
	"testing"

	"github.com/krelinga/go-deep"
	"github.com/krelinga/go-deep/exam"
	"github.com/krelinga/go-deep/match"
)

func TestMatchVals(t *testing.T) {
	// TODO: add more tests with fake matchers where I can control the output?
	tests := []struct {
		name       string
		inVal      match.Vals
		logMatcher match.Matcher
	}{
		{
			name:  "real matcher matches",
			inVal: match.NewVals1(42),
			logMatcher: matchLogOk(),
		},
		{
			name:  "real matcher does not match",
			inVal: match.NewVals1(43),
			logMatcher: matchLogError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := deep.NewEnv()
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				exam.MatchVals(e, env, tt.inVal, match.Equal(42))
			})
			result := tt.logMatcher.Match(env, match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("MatchVals() log = %s", deep.Format(env, log))
			}
		})
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name       string
		in         int
		logMatcher match.Matcher
	}{
		{
			name: "real matcher matches",
			in:   42,
			logMatcher: matchLogOk(),
		},
		{
			name: "real matcher does not match",
			in:   43,
			logMatcher: matchLogError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := deep.NewEnv()
			h := exam.Harness{}
			log := h.Run(func(e exam.E) {
				exam.Match(e, env, tt.in, match.Equal(42))
			})
			result := tt.logMatcher.Match(env, match.NewVals1(log))
			if !result.Matched() {
				t.Errorf("Match() log = %s", deep.Format(env, log))
			}
		})
	}
}
