package automapper

import (
	"github.com/pkg/errors"
	"reflect"
)

type NameMapper struct {
	template     Template
	profile      map[string]string
	transformers map[string]TransformHandler
	errs         []error
}

func NewNameMapper(template Template, profile map[string]string) (*NameMapper, error) {
	nm := &NameMapper{
		template:     template,
		profile:      profile,
		transformers: make(map[string]TransformHandler),
		errs:         make([]error, 0),
	}
	err := nm.SetProfile(profile)
	if err != nil {
		return nil, err
	}
	return nm, nil
}

func (nm *NameMapper) SetProfile(profile map[string]string) error {
	for dstFieldName, srcFieldName := range profile {
		if nm.template.indexOfDstFieldName(dstFieldName) < 0 {
			return errors.WithMessage(DoesNotExist, nm.template.dst.String()+"."+dstFieldName)
		}
		if nm.template.indexOfSrcFieldName(srcFieldName) < 0 {
			return errors.WithMessage(DoesNotExist, nm.template.src.String()+"."+srcFieldName)
		}
		nm.transformers[dstFieldName] = FromField(srcFieldName)
	}
	return nil
}

func (nm *NameMapper) mapping(srcValue, dstValue reflect.Value) error {
	srcMap, dstMap := NewFieldMapViaTemplate(srcValue, dstValue, nm.template)
	return coverPanic(func() {
		dstMap.Mapping(srcMap, nm.transformers)
	})
}

func (nm *NameMapper) listMapping(srcValue, dstValue reflect.Value) error {
	dstMax := dstValue.Len()
	srcMax := srcValue.Len()
	for i := 0; i < srcMax && i < dstMax; i++ {
		if err := nm.mapping(srcValue.Index(i), dstValue.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (nm *NameMapper) Template() Template {
	return nm.template
}

func (nm *NameMapper) Mapping(src interface{}, dst interface{}) (err error) {
	if len(nm.errs) > 0 {
		return invalidMapperWithErrors(nm.errs...)
	}
	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)
	if err = validateMappingParameter(srcValue, dstValue, nm.template); err != nil {
		return err
	}
	if isListValue(dstValue) && isListValue(srcValue) {
		return nm.listMapping(srcValue, dstValue)
	} else if isPointerValue(dstValue) && isPointerValue(srcValue) {
		return nm.mapping(srcValue.Elem(), dstValue.Elem())
	} else {
		return errors.WithMessage(InvalidParameter, "the Mapping method take Ptr, Array, Slice kinds")
	}
}

func (nm *NameMapper) MakeDestination(source interface{}) (interface{}, error) {
	var src = reflect.ValueOf(source)
	var dst reflect.Value
	var err error
	if isListValue(src) {
		dst = reflect.MakeSlice(reflect.SliceOf(nm.template.dst), src.Len(), src.Cap())
		err = nm.listMapping(src, dst)
	} else {
		dst = reflect.New(nm.template.dst).Elem()
		err = nm.mapping(src, dst)
	}
	if err != nil {
		return nil, err
	}
	return dst.Interface(), nil
}

func (nm *NameMapper) Transform(name string, handler TransformHandler) Mapper {
	if nm.template.indexOfDstFieldName(name) < 0 {
		nm.errs = append(nm.errs, errors.WithMessage(DoesNotExist, nm.template.dst.String()+"."+name))
	} else if handler == nil {
		nm.errs = append(nm.errs, errors.WithMessage(InvalidParameter, name+" handler"))
	} else {
		nm.transformers[name] = handler
	}
	return nm
}

func (nm *NameMapper) Ignore(name string) Mapper {
	if _, ok := nm.profile[name]; ok {
		nm.transformers[name] = ignore
	}
	return nm
}

func (nm *NameMapper) Condition(name string, conditionHandler ConditionHandler) Mapper {
	if nm.template.indexOfDstFieldName(name) < 0 {
		nm.errs = append(nm.errs, errors.WithMessage(DoesNotExist, nm.template.dst.String()+"."+name))
	} else {
		nm.transformers[name] = condition(nm.profile[name], conditionHandler)
	}
	return nm
}

func (nm *NameMapper) Nested(name, target string, mapper Mapper) Mapper {
	if nm.template.indexOfDstFieldName(name) < 0 {
		nm.errs = append(nm.errs, errors.WithMessage(DoesNotExist, nm.template.dst.String()+"."+name))
	} else if nm.template.indexOfSrcFieldName(target) < 0 {
		nm.errs = append(nm.errs, errors.WithMessage(DoesNotExist, nm.template.src.String()+"."+target))
	} else if mapper == nil {
		nm.errs = append(nm.errs, errors.WithMessage(InvalidParameter, name+" mapper"))
	} else {
		nm.transformers[name] = nested(target, mapper)
	}
	return nm
}
