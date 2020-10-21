package automapper

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

type Error string

const invalid = "-invalid"
const (
	DoesNotExist       Error = "does not exist or is non-public field"
	IdentifierNotFound Error = "identifier not found"
	IsNilOrInvalid     Error = "is nil or invalid"
	IsDifferentType    Error = "is different type"
	InvalidMapper      Error = "invalid mapper"
	InvalidParameter   Error = "invalid parameter"
	UndefinedError     Error = "undefined error"
)

func (e Error) Error() string {
	return string(e)
}

func invalidMapperWithErrors(errs ...error) error {
	var result error = InvalidMapper
	for _, err := range errs {
		result = errors.WithMessage(result, err.Error())
	}
	return result
}

func coverPanic(action func()) (err error) {
	defer func() {
		except := recover()
		if except == nil {
			return
		}
		switch except.(type) {
		case error:
			err = except.(error)
			break
		default:
			err = errors.WithMessage(UndefinedError, fmt.Sprint(except))
		}
	}()
	action()
	return nil
}

func mappingDescription(description string) string {
	return "Mapping(" + description + ")"
}

func isNilOrInvalidError(srcValue, dstValue reflect.Value) error {
	description := srcValue.Type().String()
	err := false
	if !srcValue.IsValid() || srcValue.IsNil() {
		description += invalid
		err = true
	}
	description += "," + dstValue.Type().String()
	if !dstValue.IsValid() || dstValue.IsNil() {
		description += invalid
		err = true
	}
	if err {
		return errors.WithMessage(IsNilOrInvalid, mappingDescription(description))
	} else {
		return nil
	}
}

func isDifferentTypeError(srcType, dstType reflect.Type, template Template) error {
	description := ""
	err := false
	description += srcType.String()
	if !template.isSrcType(srcType) {
		description += invalid
		err = true
	}
	description += "," + dstType.String()
	if !template.isDstType(dstType) {
		description += invalid
		err = true
	}
	if err {
		return errors.WithMessage(IsDifferentType, mappingDescription(description))
	} else {
		return nil
	}
}

func validateMappingParameter(srcValue, dstValue reflect.Value, template Template) (err error) {
	err = isNilOrInvalidError(srcValue, dstValue)
	if err != nil {
		return err
	}
	err = isDifferentTypeError(srcValue.Type(), dstValue.Type(), template)
	return
}
