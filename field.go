package deep

import (
	"fmt"
	"reflect"
)

type Field interface {
	GetField(t reflect.Type) reflect.StructField
	fieldIsAClosedType()
}

type namedField string

func (nf namedField) GetField(t reflect.Type) reflect.StructField {
	field, ok := t.FieldByName(string(nf))
	if !ok {
		panic(fmt.Errorf("%w: struct %s does not have field %q", ErrWrongType, t, string(nf)))
	}
	if len(field.Index) != 1 {
		panic(fmt.Errorf("%w: field %q is not a root-level field in struct %s", ErrWrongType, string(nf), t))
	}
	if field.Anonymous {
		panic(fmt.Errorf("%w: field %q is an embedded field in struct %s", ErrWrongType, string(nf), t))
	}
	return field
}

func (nf namedField) fieldIsAClosedType() {}

func NamedField(name string) Field {
	return namedField(name)
}

func EmbedTypeField(typ reflect.Type) Field {
	return embedField{Typ: typ}
}

func EmbedField[T any]() Field {
	return EmbedTypeField(reflect.TypeFor[T]())
}

type embedField struct {
	Typ reflect.Type
}

func (ef embedField) GetField(t reflect.Type) reflect.StructField {
	field, ok := t.FieldByName(t.Name())
	if !ok {
		panic(fmt.Errorf("%w: struct %s does not have field %q", ErrWrongType, t, field.Name))
	}
	if len(field.Index) != 1 {
		panic(fmt.Errorf("%w: field %q is not a root-level field in struct %s", ErrWrongType, field.Name, t))
	}
	if !field.Anonymous || field.Type != ef.Typ {
		panic(fmt.Errorf("%w: field %q is not an embedded field of type %s in struct %s", ErrWrongType, field.Name, ef.Typ, t))
	}
	return field
}

func (ef embedField) fieldIsAClosedType() {}
