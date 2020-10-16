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

func (atm AutoMapper) mapping(srcMap, dstMap FieldMap) {
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
}

func (atm AutoMapper) Mapping(src interface{}, dst interface{}) error {
	action := "Mapping"
	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)

	err := GetNilError(action, srcValue, dstValue)
	if err != nil {
		return err
	}
	err = GetNotPointerError(action, srcValue.Type(), dstValue.Type())
	if err != nil {
		return err
	}
	err = GetDifferentTypeError(action, false, srcValue.Type().Elem(), dstValue.Type().Elem(),
		atm.srcTemplateType, atm.dstTemplateType)
	if err != nil {
		return err
	}
	atm.mapping(NewFieldMapFormPointer(src), NewFieldMapFormPointer(dst))
	return nil
}

func (atm *AutoMapper) ListMapping(src interface{}, dst interface{}) error {
	action := "ListMapping"
	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)
	err := GetNilError(action, srcValue, dstValue)
	if err != nil {
		return err
	}
	err = GetNotListError(action, srcValue.Type(), dstValue.Type())
	if err != nil {
		return err
	}
	err = GetDifferentTypeError(action, true, srcValue.Type().Elem(), dstValue.Type().Elem(),
		atm.srcTemplateType, atm.dstTemplateType)
	if err != nil {
		return err
	}
	dstMax := dstValue.Len()
	srcMax := dstValue.Len()
	max := dstMax
	if srcMax < dstMax {
		max = srcMax
	}
	for i := 0; i < max; i++ {
		srcField := srcValue.Index(i)
		dstField := dstValue.Index(i)
		atm.mapping(NewFieldMapFromValue(srcField), NewFieldMapFromValue(dstField))
	}
	return nil
}
