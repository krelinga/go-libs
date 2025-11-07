package deeptest

import "github.com/krelinga/go-deep/match"

func Equal[T any](env Env, a, b T) Result {
	env.Helper()
	m := match.Equal(b)
	return Match(env, a, m)
}

func NotEqual[T any](env Env, a, b T) Result {
	env.Helper()
	m := match.NotEqual(b)
	return Match(env, a, m)
}
