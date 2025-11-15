package exam

import (
	"github.com/krelinga/go-deep/deep"
	"github.com/krelinga/go-deep/match"
)

func Nil[Type any](e E, env deep.Env, a Type) Result {
	e.Helper()
	m := match.Nil()
	return Match(e, env, a, m)
}

func NotNil[Type any](e E, env deep.Env, a Type) Result {
	e.Helper()
	m := match.Not(match.Nil())
	return Match(e, env, a, m)
}
