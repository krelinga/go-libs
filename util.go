package deep

import (
	"fmt"
	"reflect"
	"strings"
)

func indent(s string) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = "  " + lines[i]
	}
	return strings.Join(lines, "\n")
}

func Value[T any](in T) reflect.Value {
	typ := reflect.TypeFor[T]()
	v := reflect.ValueOf(in)
	if typ.Kind() == reflect.Interface {
		newV := reflect.New(typ).Elem()
		if v.IsValid() {
			newV.Set(v)
		}
		v = newV
	}
	return v
}

func typeName(t reflect.Type) string {
	if t == reflect.TypeFor[any]() {
		return "any"
	}
	return t.String()
}

func AsSlice(v reflect.Value) []reflect.Value {
	if !v.IsValid() {
		panic(ErrInvalid)
	}
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		// Expected case
	default:
		panic(fmt.Errorf("%w: value is not a slice or array: got %s", ErrWrongType, v.Kind()))
	}
	if v.IsNil() {
		return nil
	}
	out := make([]reflect.Value, v.Len())
	for i := 0; i < v.Len(); i++ {
		out[i] = v.Index(i)
	}
	return out
}
