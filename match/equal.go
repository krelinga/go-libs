package match

import "github.com/krelinga/go-deep/deep"

func Equal[T any](want T) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := Want1[T](vals)
		if deep.Equal(env, got, want) {
			return NewResult(true, "values are equal")
		}
		return NewResult(false, "values are not equal")
	})
}

func NotEqual[T any](want T) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := Want1[T](vals)
		if !deep.Equal(env, got, want) {
			return NewResult(true, "values are not equal")
		}
		return NewResult(false, "values are equal")
	})
}
