package match

import (
	"fmt"
	"reflect"

	"github.com/krelinga/go-deep"
)

type Struct struct {
	// Matchers for individual fields.
	Fields map[deep.Field]Matcher

	// If set, makes a copy of the input struct with zero values for
	// each field named in Fields, and applies Partial to that struct.
	Partial Matcher
}

func (s Struct) Match(env deep.Env, vals Vals) Result {
	got := vals.Want1()
	if got.Kind() != reflect.Struct {
		panic(fmt.Errorf("%w: value is not a struct: got %s", ErrBadType, got.Kind()))
	}
	t := got.Type()
	for fieldI, matcher := range s.Fields {
		field := fieldI.GetField(t)
		fieldVal := got.FieldByIndex(field.Index)
		result := matcher.Match(env, Vals{fieldVal})
		if !result.Matched() {
			return NewResult(false, fmt.Sprintf("field %q does not match", field.Name))
		}
	}
	if s.Partial != nil {
		partialStruct := reflect.New(got.Type()).Elem()
		for fIdx := range t.NumField() {
			field := t.Field(fIdx)
			if len(field.Index) != 1 {
				// Skip non-root-level fields.
				continue
			}
			var fieldI deep.Field
			if field.Anonymous {
				fieldI = deep.EmbedTypeField(field.Type)
			} else {
				fieldI = deep.NamedField(field.Name)
			}
			if _, ok := s.Fields[fieldI]; ok {
				// Zero out this field.
				continue
			}
			partialStruct.FieldByIndex(field.Index).Set(got.FieldByIndex(field.Index))
		}
		result := s.Partial.Match(env, Vals{partialStruct})
		if !result.Matched() {
			return NewResult(false, "partial struct does not match")
		}
	}
	return NewResult(true, "struct matches")
}
