package deep

import "reflect"

type Field interface {
	fieldIsAClosedType()
}

type namedField string

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

func (ef embedField) fieldIsAClosedType() {}
