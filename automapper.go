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
	for tag, field := range dstMap {
		handler, ok := atm.transforms[tag]
		handlerValue := reflect.ValueOf(handler)
		if ok &&
			handler != nil &&
			handlerValue != ignoreHandlerValue &&
			handlerValue != defaultHandlerValue {
			value := reflect.ValueOf(handler(srcMap))
			if value == ignoreHandlerValue {
				continue
			}
			if value != defaultHandlerValue {
				field.Set(value)
				continue
			}
		}
		srcValue, ok := srcMap[tag]
		if ok {
			field.Set(srcValue)
		}
	}
}

func (atm AutoMapper) Mapping(src interface{}, dst interface{}) error {
	action := "Mapping"
	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)

	err := IsNilOrInvalidError(action, srcValue, dstValue)
	if err != nil {
		return err
	}
	err = IsNotPointerError(action, srcValue.Type(), dstValue.Type())
	if err != nil {
		return err
	}
	err = IsDifferentTypeError(action, false, srcValue.Type().Elem(), dstValue.Type().Elem(),
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
	err := IsNilOrInvalidError(action, srcValue, dstValue)
	if err != nil {
		return err
	}
	err = IsNotListError(action, srcValue.Type(), dstValue.Type())
	if err != nil {
		return err
	}
	err = IsDifferentTypeError(action, true, srcValue.Type().Elem(), dstValue.Type().Elem(),
		atm.srcTemplateType, atm.dstTemplateType)
	if err != nil {
		return err
	}
	dstMax := dstValue.Len()
	srcMax := srcValue.Len()
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
