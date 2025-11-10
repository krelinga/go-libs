package deeptest_test

import (
	"testing"

	"github.com/krelinga/go-deep"
	"github.com/krelinga/go-deep/deeptest"
)

func TestHarness(t *testing.T) {
	t.Run("env propagation", func(t *testing.T) {
		h := deeptest.Harness{
			DeepEnv: deep.NewEnv(testValOpt("foo")),
		}
		h.Run(func(innerT deeptest.T) {
			val := getTestVal(t, innerT.DeepEnv())
			if val != "foo" {
				t.Errorf("expected test value 'foo', got %q", val)
			}
			innerT.Run("nested", func(innerT deeptest.T) {
				val := getTestVal(t, innerT.DeepEnv())
				if val != "foo" {
					t.Errorf("expected test value 'foo', got %q", val)
				}
			})
		})
	})
}
