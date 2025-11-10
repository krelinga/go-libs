package deeptest_test

import (
	"reflect"
	"testing"

	"github.com/krelinga/go-deep"
)

type testValKey struct{}

func getTestVal(t *testing.T, env deep.Env) string {
	t.Helper()
	val, ok := env.Get(reflect.TypeOf(""), testValKey{})
	if !ok {
		t.Fatal("test value not found in env")
	}
	return val.(string)
}

func testValOpt(val string) deep.Opt {
	return deep.OptFunc(func(env deep.EnvSetter) {
		env.Set(reflect.TypeOf(""), testValKey{}, val)
	})
}