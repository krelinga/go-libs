package exam

import (
	"github.com/krelinga/go-libs/deep"
	"github.com/krelinga/go-libs/match"
)

func PanicWith(e E, env deep.Env, m match.Matcher, f func()) (res Result) {
	e.Helper()
	defer func() {
		if r := recover(); r != nil {
			v := match.NewVals1(r)
			matchResult := m.Match(env, v)
			if !matchResult.Matched() {
				res = NewResult(false, e)
				e.Errorf("function panicked, but panic did not match: %s", matchResult)
			} else {
				res = NewResult(true, e)
			}
		} else {
			e.Errorf("expected panic, but function completed without panicking")
			res = NewResult(false, e)
		}
	}()
	f()
	return
}

func Panic(e E, env deep.Env, f func()) Result {
	return PanicWith(e, env, match.True(), f)
}