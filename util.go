package automapper

import "reflect"

func CompareType(obj interface{}, t reflect.Type) bool {
	return reflect.TypeOf(obj).Elem() == t
}

func IsNil(obj interface{}) bool {
	return reflect.ValueOf(obj).IsNil()
}

func IsPointer(obj interface{}) bool {
	return reflect.ValueOf(obj).Kind() == reflect.Ptr
}
