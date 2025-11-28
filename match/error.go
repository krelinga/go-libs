package match

import (
	"errors"

	"github.com/krelinga/go-libs/deep"
)

func ErrorIs(target error) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := Want1[error](vals)
		if !errors.Is(got, target) {
			return NewResult(false, "error is not target")
		}
		return NewResult(true, "error is target")
	})
}
