package automapper

import "reflect"

// private
var __ignore = func() string { return TAG }

func IGNORE() interface{} { return __ignore }

func ignore(SrcMap FieldMap) interface{} { return __ignore }

func condition(name string, condition ConditionHandler) TransformHandler {
	return func(SrcMap FieldMap) interface{} {
		if condition(SrcMap) {
			return SrcMap.Interface(name)
		} else {
			return ignore(SrcMap)
		}
	}
}

func nested(target string, mapper Mapper) TransformHandler {
	return func(SrcMap FieldMap) interface{} {
		value, err := mapper.MakeDestination(SrcMap.Interface(target))
		if err != nil {
			panic(err)
		}
		return value
	}
}

func isUppercaseCharacter(c uint8) bool {
	return c >= 65 && c <= 90
}

func isPublicField(field reflect.StructField) bool {
	if len(field.Name) < 1 {
		return false
	}
	return isUppercaseCharacter(field.Name[0])
}

func isIgnore(value reflect.Value) bool {
	if !value.IsValid() {
		return false
	} else if value.Kind() != reflect.Func {
		return false
	} else if value.IsNil() {
		return false
	} else {
		ignValue := reflect.ValueOf(__ignore)
		return ignValue == value
	}
}

func isListValue(value reflect.Value) bool {
	kind := value.Kind()
	return kind == reflect.Array || kind == reflect.Slice
}

func isPointerValue(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr
}

func isListType(value reflect.Type) bool {
	kind := value.Kind()
	return kind == reflect.Array || kind == reflect.Slice
}

func isPointerType(value reflect.Type) bool {
	return value.Kind() == reflect.Ptr
}
