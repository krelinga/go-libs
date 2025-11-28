package match

import (
	"fmt"
	"reflect"

	"github.com/krelinga/go-libs/deep"
	"github.com/krelinga/go-libs/zero"
)

type Vals []reflect.Value

func (v Vals) Check(wantLen int) {
	if len(v) != wantLen {
		panic(fmt.Errorf("%w: got %d, want %d", ErrBadVals, len(v), wantLen))
	}
	for i, val := range v {
		if !val.IsValid() {
			panic(fmt.Errorf("%w: val %d is not valid", ErrBadVals, i))
		}
		if !val.CanInterface() {
			panic(fmt.Errorf("%w: val %d cannot interface", ErrBadVals, i))
		}
	}
}

func (v Vals) Want1() reflect.Value {
	v.Check(1)
	return v[0]
}

func (v Vals) Want2() (reflect.Value, reflect.Value) {
	v.Check(2)
	return v[0], v[1]
}

func (v Vals) Want3() (reflect.Value, reflect.Value, reflect.Value) {
	v.Check(3)
	return v[0], v[1], v[2]
}

func NewValsAny(vals ...any) Vals {
	v := make(Vals, len(vals))
	for i, val := range vals {
		v[i] = reflect.ValueOf(val)
	}
	return v
}

func NewVals1[T any](v T) Vals {
	return Vals{deep.Value(v)}
}

func NewVals2[T1, T2 any](v1 T1, v2 T2) Vals {
	return Vals{deep.Value(v1), deep.Value(v2)}
}

func NewVals3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3) Vals {
	return Vals{deep.Value(v1), deep.Value(v2), deep.Value(v3)}
}

func Want1[T any](vals Vals) T {
	vals.Check(1)
	return getType[T](vals[0])
}

func Want2[T1, T2 any](vals Vals) (T1, T2) {
	vals.Check(2)
	return getType[T1](vals[0]), getType[T2](vals[1])
}

func Want3[T1, T2, T3 any](vals Vals) (T1, T2, T3) {
	vals.Check(3)
	return getType[T1](vals[0]), getType[T2](vals[1]), getType[T3](vals[2])
}

func getType[T any](v reflect.Value) T {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		if !v.Type().Implements(t) {
			panic(fmt.Errorf("%w: value of type %s does not implement %s", ErrBadVals, v.Type(), t))
		}
		if v.IsNil() {
			return zero.For[T]()
		}
	}
	return v.Interface().(T)
}