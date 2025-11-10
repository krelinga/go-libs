package deeptest

import (
	"github.com/krelinga/go-deep"
	"github.com/krelinga/go-deep/match"
)

func Equal[Type any](t T, env deep.Env, a, b Type) Result {
	t.Helper()
	m := match.Equal(b)
	return Match(t, env, a, m)
}

func NotEqual[Type any](t T, env deep.Env, a, b Type) Result {
	t.Helper()
	m := match.NotEqual(b)
	return Match(t, env, a, m)
}
