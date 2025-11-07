package match

import "github.com/krelinga/go-deep"

func LessThan[T any](limit T) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := Want1[T](vals)
		if deep.Order(env, got, limit) < 0 {
			return NewResult(true, "value is less than limit")
		}
		return NewResult(false, "value is not less than limit")
	})
}

func LessThanOrEqual[T any](limit T) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := Want1[T](vals)
		if deep.Order(env, got, limit) <= 0 {
			return NewResult(true, "value is less than or equal to limit")
		}
		return NewResult(false, "value is greater than limit")
	})
}

func GreaterThan[T any](limit T) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := Want1[T](vals)
		if deep.Order(env, got, limit) > 0 {
			return NewResult(true, "value is greater than limit")
		}
		return NewResult(false, "value is not greater than limit")
	})
}

func GreaterThanOrEqual[T any](limit T) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := Want1[T](vals)
		if deep.Order(env, got, limit) >= 0 {
			return NewResult(true, "value is greater than or equal to limit")
		}
		return NewResult(false, "value is less than limit")
	})
}
