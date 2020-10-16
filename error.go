package automapper

import (
	"errors"
	"reflect"
)

type DoesNotExist struct {
	fieldName string
}

func (e *DoesNotExist) Error() string {
	return e.fieldName + " does not exist"
}

type autoMapperValueError struct {
	dstValue reflect.Value
	srcValue reflect.Value
	action   string
}

type autoMapperTypeError struct {
	dstType reflect.Type
	srcType reflect.Type
	action  string
}

func GetNilError(action string, srcValue, dstValue reflect.Value) error {
	err := ""
	isErr := false
	if !srcValue.IsValid() || srcValue.IsNil() {
		err += "is nil,"
		isErr = true
	} else {
		err += "src,"
	}
	if !dstValue.IsValid() || dstValue.IsNil() {
		err += "is nil"
		isErr = true
	} else {
		err += "dst"
	}
	err = "automapper: " + action + "(" + err + ")"
	if isErr {
		return errors.New(err)
	} else {
		return nil
	}
}

func GetNotPointerError(action string, srcType, dstType reflect.Type) error {
	err := ""
	isErr := false
	err += srcType.String()
	if srcType.Kind() != reflect.Ptr {
		err += " isn't pointer"
		isErr = true
	}
	err += ", " + dstType.String()
	if dstType.Kind() != reflect.Ptr {
		err += " isn't pointer"
		isErr = true
	}
	err = "automapper: " + action + "(" + err + ")"
	if isErr {
		return errors.New(err)
	} else {
		return nil
	}
}

func GetNotListError(action string, srcType, dstType reflect.Type) error {
	err := ""
	isErr := false
	err += srcType.String()
	if srcType.Kind() != reflect.Slice {
		err += " isn't list"
		isErr = true
	}
	err += ", " + dstType.String()
	if dstType.Kind() != reflect.Slice {
		err += " isn't list"
		isErr = true
	}
	err = "automapper: " + action + "(" + err + ")"
	if isErr {
		return errors.New(err)
	} else {
		return nil
	}
}

func GetDifferentTypeError(action string, many bool, srcType, dstType, srcTmpl, dstTmpl reflect.Type) error {
	err := ""
	isErr := false
	if srcType != srcTmpl {
		if many {
			err += "isn't []"
		} else {
			err += "isn't "
		}
		isErr = true
	}
	err += srcTmpl.String() + ","
	if dstType != dstTmpl {
		if many {
			err += "isn't []"
		} else {
			err += "isn't "
		}
		isErr = true
	}
	err += dstTmpl.String()
	err = "automapper: " + action + "(" + err + ")"
	if isErr {
		return errors.New(err)
	} else {
		return nil
	}
}
