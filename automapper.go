package automapper

const TAG = "automapper"

type ConditionHandler func(SrcMap FieldMap) bool

type Mapper interface {
	GetTemplate
	Mapping
	MakeDestination
	Transform
	Ignore
	Condition
	Nested
}
type GetTemplate interface {
	Template() Template
}

type Mapping interface {
	Mapping(src, dst interface{}) error
}

type Transform interface {
	Transform(identifier string, handler TransformHandler) Mapper
}

type Ignore interface {
	Ignore(identifier string) Mapper
}

type Condition interface {
	Condition(identifier string, conditionHandler ConditionHandler) Mapper
}

type Nested interface {
	Nested(name, target string, mapper Mapper) Mapper
}

type MakeDestination interface {
	MakeDestination(src interface{}) (interface{}, error)
}

func FromField(target string) TransformHandler {
	return func(SrcMap FieldMap) interface{} {
		return SrcMap.elements[target].Interface()
	}
}
