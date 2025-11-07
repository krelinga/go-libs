package match

import "github.com/krelinga/go-deep"

type Matcher interface {
	Match(deep.Env, Vals) Result
}

type Func func(deep.Env, Vals) Result

func (f Func) Match(env deep.Env, vals Vals) Result {
	return f(env, vals)
}