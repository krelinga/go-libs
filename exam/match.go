package exam

import (
	"github.com/krelinga/go-deep/deep"
	"github.com/krelinga/go-deep/match"
)

func Match[Type any](e E, env deep.Env, got Type, m match.Matcher) Result {
	e.Helper()
	vals := match.NewVals1(got)
	return MatchVals(e, env, vals, m)
}

func MatchVals(e E, env deep.Env, vals match.Vals, m match.Matcher) Result {
	e.Helper()
	r := m.Match(env, vals)
	if !r.Matched() {
		e.Errorf("Match failed:\n%s", r)
	}
	return NewResult(r.Matched(), e)
}
