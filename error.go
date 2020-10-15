package automapper

import "reflect"

type DoesNotExist struct {
	fieldName string
}

func (e *DoesNotExist) Error() string {
	return e.fieldName + " does not exist"
}

type InvalidMappingError struct {
	DstType reflect.Type
	SrcType reflect.Type
}

func (e *InvalidMappingError) Error() string {
	er := ""
	if e.DstType == nil || e.SrcType == nil {
		if e.SrcType == nil {
			er += "nil,"
		} else {
			er += "src,"
		}
		if e.DstType == nil {
			er += "nil"
		} else {
			er += "dst"
		}
	} else if e.SrcType.Kind() != reflect.Ptr || e.DstType.Kind() != reflect.Ptr {
		if e.SrcType.Kind() != reflect.Ptr {
			er += "non-pointer :" + e.SrcType.String() + ","
		} else {
			er += e.SrcType.String() + ","
		}
		if e.DstType.Kind() != reflect.Ptr {
			er += "non-pointer :" + e.DstType.String()
		} else {
			er += e.DstType.String()
		}
	} else {
		return "automapper: Transform invalid-type"
	}
	return "automapper: Transform(" + er + ")"
}
