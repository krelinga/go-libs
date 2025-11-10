package deeptest

import (
	"github.com/krelinga/go-deep"
	"github.com/krelinga/go-deep/match"
)

func Match[Type any](t T, env deep.Env, got Type, m match.Matcher) Result {
	t.Helper()
	vals := match.NewVals1(got)
	return MatchVals(t, env, vals, m)
}

func MatchVals(t T, env deep.Env, vals match.Vals, m match.Matcher) Result {
	t.Helper()
	r := m.Match(env, vals)
	if !r.Matched() {
		t.Errorf("Match failed:\n%s", r)
	}
	return &resultImpl{
		ok:  r.Matched(),
		t: t,
		env: env,
	}
}
