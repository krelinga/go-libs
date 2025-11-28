package match

import (
	"fmt"
	"reflect"

	"github.com/krelinga/go-libs/deep"
)

func As[T any](m Matcher) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		v := vals.Want1()
		if v.Type().Kind() == reflect.Interface && !v.IsNil() {
			v = v.Elem()
		}
		t := reflect.TypeFor[T]()
		if !v.Type().ConvertibleTo(t) {
			return NewResult(false, fmt.Sprintf("value is not convertible to type %s: got %s", t, v.Type()))
		}
		return m.Match(env, Vals{v.Convert(t)})
	})
}