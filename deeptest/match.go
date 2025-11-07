package deeptest

import "github.com/krelinga/go-deep/match"

func Match[T any](env Env, got T, m match.Matcher) Result {
	vals := match.NewVals1(got)
	r := m.Match(env, vals)
	if !r.Matched() {
		env.Errorf("Match failed:\n%s", r)
	}
	return &resultImpl{
		ok:  r.Matched(),
		env: env,
	}
}
