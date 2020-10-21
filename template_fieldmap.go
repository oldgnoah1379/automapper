package automapper

import (
	"reflect"
)

func NewFieldMapViaTemplate(src, dst reflect.Value, t Template) (SrcMap FieldMap, DstMap FieldMap) {
	SrcMap = FieldMap{elements: make(map[string]reflect.Value)}
	DstMap = FieldMap{elements: make(map[string]reflect.Value)}

	for name, idx := range t.srcNames {
		element := src.Field(idx)
		if element.IsValid() {
			SrcMap.elements[name] = element
		}
	}
	for name, idx := range t.dstNames {
		element := dst.Field(idx)
		if element.IsValid() {
			DstMap.elements[name] = element
		}
	}
	return
}
