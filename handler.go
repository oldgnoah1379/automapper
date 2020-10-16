package automapper

import "reflect"

type TransformHandler func(SrcMap FieldMap) interface{}

var Ignore TransformHandler = func(SrcMap FieldMap) interface{} {
	return TAG
}
var Default TransformHandler = func(SrcMap FieldMap) interface{} {
	return TAG + "Default"
}
var ignoreHandlerValue = reflect.ValueOf(Ignore)
var defaultHandlerValue = reflect.ValueOf(Default)

func Condition(conditionCallback func(SrcMap FieldMap) bool, transform TransformHandler) TransformHandler {
	return func(SrcMap FieldMap) interface{} {
		if conditionCallback(SrcMap) {
			if transform != nil {
				return transform(SrcMap)
			}
			return Default
		} else {
			return Ignore
		}
	}
}
