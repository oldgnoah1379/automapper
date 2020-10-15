package automapper

type TransformHandler func(SrcMap FieldMap) interface{}

var Ignore TransformHandler = func(SrcMap FieldMap) interface{} {
	return TAG
}
