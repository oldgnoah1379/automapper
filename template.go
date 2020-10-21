package automapper

import (
	"reflect"
)

type Template struct {
	dst      reflect.Type
	src      reflect.Type
	dstNames map[string]int
	dstTags  map[string]int
	srcNames map[string]int
	srcTags  map[string]int
}

func NewTemplate(src, dst interface{}) Template {
	t := Template{
		src:      reflect.TypeOf(src),
		srcNames: make(map[string]int),
		srcTags:  make(map[string]int),
		dst:      reflect.TypeOf(dst),
		dstTags:  make(map[string]int),
		dstNames: make(map[string]int),
	}
	if t.src.Kind() == reflect.Ptr {
		t.src = t.src.Elem()
	}
	if t.dst.Kind() == reflect.Ptr {
		t.dst = t.dst.Elem()
	}
	numFields := t.src.NumField()
	for i := 0; i < numFields; i++ {
		field := t.src.Field(i)
		if !isPublicField(field) {
			continue
		}
		t.srcNames[field.Name] = field.Index[0]
		if tag := field.Tag.Get(TAG); len(tag) > 0 {
			t.srcTags[tag] = field.Index[0]
		}
	}
	numFields = t.dst.NumField()
	for i := 0; i < numFields; i++ {
		field := t.dst.Field(i)
		if !isPublicField(field) {
			continue
		}
		t.dstNames[field.Name] = field.Index[0]
		if tag := field.Tag.Get(TAG); len(tag) > 0 {
			t.dstTags[tag] = field.Index[0]
		}
	}
	return t
}

func (t Template) ProfileSameName() map[string]string {
	profile := make(map[string]string)
	for name, _ := range t.dstNames {
		if _, ok := t.srcNames[name]; ok {
			profile[name] = name
		}
	}
	return profile
}

func (t Template) ProfileSameTag() map[string]string {
	profile := make(map[string]string)
	for tag, dstIdx := range t.dstTags {
		if srcIdx, ok := t.srcTags[tag]; ok {
			profile[t.dst.Field(dstIdx).Name] = t.src.Field(srcIdx).Name
		}
	}
	return profile
}

func (t Template) indexOfSrcFieldTag(tag string) int {
	idx, ok := t.srcTags[tag]
	if !ok {
		return -1
	}
	return idx
}

func (t Template) indexOfDstFieldTag(tag string) int {
	idx, ok := t.dstTags[tag]
	if !ok {
		return -1
	}
	return idx
}

func (t Template) indexOfSrcFieldName(name string) int {
	idx, ok := t.srcNames[name]
	if !ok {
		return -1
	}
	return idx
}

func (t Template) indexOfDstFieldName(name string) int {
	idx, ok := t.dstNames[name]
	if !ok {
		return -1
	}
	return idx
}

func (t Template) isSrcType(compareType reflect.Type) bool {
	return t.isType(compareType, true)
}

func (t Template) isDstType(compareType reflect.Type) bool {
	return t.isType(compareType, false)
}

func (t Template) isType(compareType reflect.Type, isSrc bool) bool {
	template := t.dst
	if isSrc {
		template = t.src
	}
	if isPointerType(compareType) || isListType(compareType) {
		return compareType.Elem() == template
	} else {
		return compareType == template
	}
}
