package automapper

import "reflect"

type FieldMap map[string]reflect.Value

func NewFieldMap(obj interface{}) FieldMap {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()
	result := make(map[string]reflect.Value)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		result[field.Tag.Get(TAG)] = v.Field(i)
	}
	return result
}

func (f FieldMap) Field(fieldName string) reflect.Value {
	v, ok := f[fieldName]
	if !ok {
		exist := &DoesNotExist{fieldName: fieldName}
		panic(exist.Error())
	}
	var r = v.Interface()
	return reflect.ValueOf(r)
}

func (f FieldMap) String(fieldName string) string {
	return f.Field(fieldName).String()
}

func (f FieldMap) Int(fieldName string) int {
	return int(f.Field(fieldName).Int())
}

func (f FieldMap) Float(fieldName string) float64 {
	return f.Field(fieldName).Float()
}

func (f FieldMap) Bytes(fieldName string) []byte {
	return f.Field(fieldName).Bytes()
}

func (f FieldMap) Bool(fieldName string) bool {
	return f.Field(fieldName).Bool()
}

func (f FieldMap) Interface(fieldName string) interface{} {
	return f.Field(fieldName).Interface()
}
