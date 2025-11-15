package exam

import (
	"github.com/krelinga/go-deep/deep"
	"github.com/krelinga/go-deep/match"
)

func Equal[Type any](e E, env deep.Env, a, b Type) Result {
	e.Helper()
	m := match.Equal(b)
	return Match(e, env, a, m)
}

func NotEqual[Type any](e E, env deep.Env, a, b Type) Result {
	e.Helper()
	m := match.NotEqual(b)
	return Match(e, env, a, m)
}
