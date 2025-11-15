package match

import (
	"fmt"
	"reflect"

	"github.com/krelinga/go-deep/deep"
)

func Pointer(m Matcher) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := vals.Want1()
		if got.Kind() != reflect.Ptr {
			panic(fmt.Errorf("%w: value is not a pointer: got %s", ErrBadType, got.Kind()))
		}
		if got.IsNil() {
			return NewResult(false, "pointer is nil")
		}
		elem := got.Elem()
		return m.Match(env, Vals{elem})
	})
}
