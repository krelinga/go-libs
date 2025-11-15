package match

import (
	"fmt"
	"reflect"

	"github.com/krelinga/go-deep/deep"
)

func Nil() Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := vals.Want1()
		switch got.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			// Expected case
		default:
			panic(fmt.Errorf("%w: value does not support the Nil operation: %s", ErrBadType, got.Kind()))
		}
		if got.IsNil() {
			return NewResult(true, "value is nil")
		}
		return NewResult(false, "value is not nil")
	})
}
