package match

import "github.com/krelinga/go-deep/deep"

func AllOf(matchers ...Matcher) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		for _, matcher := range matchers {
			result := matcher.Match(env, vals)
			if !result.Matched() {
				return NewResult(false, "not all matchers matched: "+result.String())
			}
		}
		return NewResult(true, "all matchers matched")
	})
}

func AnyOf(matchers ...Matcher) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		for _, matcher := range matchers {
			result := matcher.Match(env, vals)
			if result.Matched() {
				return NewResult(true, "at least one matcher matched: "+result.String())
			}
		}
		return NewResult(false, "no matchers matched")
	})
}

func Not(matcher Matcher) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		result := matcher.Match(env, vals)
		if result.Matched() {
			return NewResult(false, "matcher matched: "+result.String())
		}
		return NewResult(true, "matcher did not match: "+result.String())
	})
}

func NoneOf(matchers ...Matcher) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		for _, matcher := range matchers {
			result := matcher.Match(env, vals)
			if result.Matched() {
				return NewResult(false, "at least one matcher matched: "+result.String())
			}
		}
		return NewResult(true, "no matchers matched")
	})
}

func True() Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		return NewResult(true, "always true")
	})
}

func False() Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		return NewResult(false, "always false")
	})
}
