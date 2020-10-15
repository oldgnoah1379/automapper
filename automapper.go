package automapper

import (
	"reflect"
)

const TAG = "automapper"

type AutoMapper struct {
	transforms      map[string]TransformHandler
	dstTemplateType reflect.Type
	srcTemplateType reflect.Type
}

func (atm *AutoMapper) Mapping(src interface{}, dst interface{}) error {
	err := &InvalidMappingError{SrcType: reflect.TypeOf(src), DstType: reflect.TypeOf(dst)}
	if !CompareType(src, atm.srcTemplateType) ||
		!CompareType(dst, atm.dstTemplateType) ||
		!IsPointer(src) ||
		!IsPointer(dst) ||
		IsNil(src) ||
		IsNil(dst) {
		return err
	}
	srcMap := NewFieldMap(src)
	dstMap := NewFieldMap(dst)
	ign := reflect.ValueOf(Ignore)
	for tag, field := range dstMap {
		handler, ok := atm.transforms[tag]
		if ok {
			if reflect.ValueOf(handler) != ign {
				value := handler(srcMap)
				field.Set(reflect.ValueOf(value))
			}
			continue
		}
		srcValue, ok := srcMap[tag]
		if ok {
			field.Set(srcValue)
			continue
		}
	}
	return nil
}
