package match

import (
	"fmt"
	"reflect"

	"github.com/krelinga/go-deep"
)

func Len(m Matcher) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := vals.Want1()
		switch got.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			// Expected case
		default:
			panic(fmt.Errorf("%w: value does not support the Len operation: %s", ErrBadType, got.Kind()))
		}
		length := got.Len()
		return m.Match(env, NewVals1(length))
	})
}
