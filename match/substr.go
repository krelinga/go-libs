package match

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/krelinga/go-deep/deep"
)

func Substring(sub string) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := vals.Want1()
		if got.Kind() != reflect.String {
			panic(fmt.Errorf("%w: value is not a string: got %s", ErrBadType, got.Kind()))
		}
		gotString := got.String()
		if strings.Contains(gotString, sub) {
			return NewResult(true, fmt.Sprintf("string contains substring %q", sub))
		}
		return NewResult(false, fmt.Sprintf("string does not contain substring %q", sub))
	})
}
