package automapper

import (
	"github.com/pkg/errors"
	"reflect"
)

var (
	DoesNotExist    = errors.New("does not exist")
	IsNilOrInvalid  = errors.New("is nil or invalid")
	IsNotPointer    = errors.New("isn't pointer")
	IsNotList       = errors.New("isn't list")
	IsDifferentType = errors.New("is different type")
	Invalid         = "-invalid"
)

func GetErrorDescription(action, description string) string {
	return "automapper." + action + "(" + description + ")"
}

func IsNilOrInvalidError(action string, srcValue, dstValue reflect.Value) error {
	description := srcValue.String()
	err := false
	if !srcValue.IsValid() || srcValue.IsNil() {
		description += Invalid
		err = true
	}
	description += "," + dstValue.String()
	if !dstValue.IsValid() || dstValue.IsNil() {
		description += Invalid
		err = true
	}
	if err {
		return errors.WithMessage(IsNilOrInvalid, GetErrorDescription(action, description))
	} else {
		return nil
	}
}

func IsNotPointerError(action string, srcType, dstType reflect.Type) error {
	description := srcType.String()
	err := false
	if srcType.Kind() != reflect.Ptr {
		description += Invalid
		err = true
	}
	description += "," + dstType.String()
	if dstType.Kind() != reflect.Ptr {
		description += Invalid
		err = true
	}
	if err {
		return errors.WithMessage(IsNotPointer, GetErrorDescription(action, description))
	} else {
		return nil
	}
}

func IsNotListError(action string, srcType, dstType reflect.Type) error {
	description := srcType.String()
	err := false
	srcKind := srcType.Kind()
	if srcKind != reflect.Slice && srcKind != reflect.Array {
		description += Invalid
		err = true
	}
	description += "," + dstType.String()
	dstKind := dstType.Kind()
	if dstKind != reflect.Slice && dstKind != reflect.Array {
		description += Invalid
		err = true
	}
	if err {
		return errors.WithMessage(IsNotList, GetErrorDescription(action, description))
	} else {
		return nil
	}
}

func IsDifferentTypeError(action string, many bool, srcType, dstType, srcTmpl, dstTmpl reflect.Type) error {
	description := ""
	err := false
	if many {
		description += "[]"
	}
	description += srcType.String()
	if srcType != srcTmpl {
		description += Invalid
		err = true
	}
	description += ","
	if many {
		description += "[]"
	}
	description += dstType.String()
	if dstType != dstTmpl {
		description += Invalid
		err = true
	}
	if err {
		return errors.WithMessage(IsDifferentType, GetErrorDescription(action, description))
	} else {
		return nil
	}
}
