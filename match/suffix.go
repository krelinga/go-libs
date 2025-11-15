package match

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/krelinga/go-deep/deep"
)

func Suffix[S Sequence](seq S) Matcher {
	return Func(func(env deep.Env, vals Vals) Result {
		got := vals.Want1()
		switch k := reflect.TypeFor[S]().Kind(); k {
		case reflect.String:
			if got.Kind() != reflect.String {
				panic(fmt.Errorf("%w: sequence is a string but value is not: got %s", ErrBadType, got.Kind()))
			}
			gotString := got.String()
			suffix := reflect.ValueOf(seq).String()
			if strings.HasSuffix(gotString, suffix) {
				return NewResult(true, fmt.Sprintf("string has suffix %q", suffix))
			}
			return NewResult(false, fmt.Sprintf("string does not have suffix %q", suffix))
		case reflect.Slice:
			gotSlice := deep.AsSlice(got)
			matchers := reflect.ValueOf(seq).Interface().([]Matcher)
			if len(gotSlice) < len(matchers) {
				return NewResult(false, fmt.Sprintf("slice length %d is less than suffix length %d", len(gotSlice), len(matchers)))
			}
			start := len(gotSlice) - len(matchers)
			for i, m := range matchers {
				res := m.Match(env, Vals{gotSlice[start+i]})
				if !res.Matched() {
					return NewResult(false, fmt.Sprintf("slice element %d does not match", start+i))
				}
			}
			return NewResult(true, fmt.Sprintf("slice matches suffix of length %d", len(matchers)))
		default:
			panic(fmt.Errorf("%w: sequence is not a string or slice: got %s", ErrInternal, k))
		}
	})
}
