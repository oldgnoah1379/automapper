package automapper

import (
	"github.com/pkg/errors"
	"reflect"
)

type TransformHandler func(SrcMap FieldMap) interface{}
type FieldMap struct {
	elements map[string]reflect.Value
}

func (fm FieldMap) Mapping(SrcMap FieldMap, transformers map[string]TransformHandler) {
	for name, field := range fm.elements {
		handler, ok := transformers[name]
		if ok {
			value := reflect.ValueOf(handler(SrcMap))
			if !isIgnore(value) {
				field.Set(value)
			}
		}
	}
}

func (fm FieldMap) Field(identifier string) reflect.Value {
	v, ok := fm.elements[identifier]
	if !ok {
		panic(errors.WithMessage(IdentifierNotFound, identifier))
	}
	return reflect.ValueOf(v.Interface())
}

func (fm FieldMap) String(identifier string) string {
	return fm.Field(identifier).String()
}

func (fm FieldMap) Int64(identifier string) int64 {
	return fm.Field(identifier).Int()
}

func (fm FieldMap) Int(identifier string) int {
	return int(fm.Field(identifier).Int())
}

func (fm FieldMap) Float(identifier string) float64 {
	return fm.Field(identifier).Float()
}

func (fm FieldMap) Complex(identifier string) complex128 {
	return fm.Field(identifier).Complex()
}

func (fm FieldMap) Bytes(identifier string) []byte {
	return fm.Field(identifier).Bytes()
}

func (fm FieldMap) Bool(identifier string) bool {
	return fm.Field(identifier).Bool()
}

func (fm FieldMap) Interface(identifier string) interface{} {
	return fm.Field(identifier).Interface()
}
