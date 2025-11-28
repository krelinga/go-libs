package match_test

import (
	"testing"

	"github.com/krelinga/go-libs/match"
	"github.com/krelinga/go-libs/zero"
)

type myInterface interface {
	DoSomething()
}

type myStruct struct{}

func (m myStruct) DoSomething() {}

func runPanic(t *testing.T, name string, f func(t *testing.T)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				// recovered from panic
			} else {
				t.Errorf("expected panic, but function completed without panicking")
			}
		}()
		f(t)
	})
}

func runNoPanic(t *testing.T, name string, f func(t *testing.T)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("expected no panic, but function panicked: %v", r)
			}
		}()
		f(t)
	})
}

func TestWant(t *testing.T) {
	t.Run("Interface Type", func(t *testing.T) {
		runNoPanic(t, "myStruct implements myInterface non-nil", func(t *testing.T) {
			v := match.NewVals1[myInterface](myStruct{})
			match.Want1[myInterface](v) // Should not panic
		})
		runNoPanic(t, "myStruct implements myInterface nil", func(t *testing.T) {
			v := match.NewVals1(zero.For[myInterface]())
			match.Want1[myInterface](v) // Should not panic
		})
		runPanic(t, "int does not implement myInterface", func(t *testing.T) {
			v := match.NewVals1(5)
			match.Want1[myInterface](v) // Should panic
		})
	})
	t.Run("Concrete Type", func(t *testing.T) {
		runNoPanic(t, "Right Type", func(t *testing.T) {
			v := match.NewVals1(5)
			val := match.Want1[int](v)
			if val != 5 {
				t.Errorf("Want1[int] returned %v, want 5", val)
			}
		})
		runPanic(t, "Wrong Type", func(t *testing.T) {
			v := match.NewVals1("string")
			match.Want1[int](v) // Should panic
		})
	})
	t.Run("Wrong Length", func(t *testing.T) {
		runPanic(t, "Expecting 2 values but got 1", func(t *testing.T) {
			v := match.NewVals1(5)
			match.Want2[int, int](v) // Should panic
		})
		runPanic(t, "Expecting 1 value but got 2", func(t *testing.T) {
			v := match.NewVals2(5, 10)
			match.Want1[int](v) // Should panic
		})
	})
}
