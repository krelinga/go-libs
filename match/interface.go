package match

import (
	"fmt"
	"reflect"

	"github.com/krelinga/go-deep/deep"
)

func Interface(m Matcher) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := vals.Want1()
		if got.Kind() != reflect.Interface {
			panic(fmt.Errorf("%w: value is not an interface: got %s", ErrBadType, got.Kind()))
		}
		if got.IsNil() {
			return NewResult(false, "interface is nil")
		}
		elem := got.Elem()
		return m.Match(env, Vals{elem})
	})
}
