package match

import "github.com/krelinga/go-deep/deep"

func Zero() Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		val := vals.Want1()
		if val.IsZero() {
			return NewResult(true, "Is zero value")
		}
		return NewResult(false, "Is not zero value")
	})
}
