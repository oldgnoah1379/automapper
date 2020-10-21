# Automapper for Golang
## Introduction
The Automapper is the package for Go programs mapping values in different structs

The import path for the package is `github.com/gnoah1379/automapper`\
To install it, run:
```
go get github.com/gnoah1379/automapper
```
## Features
* Support mapping struct-to-struct
* Support mapping slice/array-to-slice/array
* Support Custom transform values
* Support Nested mapping (including slice/array)
## Mapper
Mapper is an interface
```Go
type Mapper interface {
	Type() string
	Mapping(src, dst interface{}) error
	MakeDestination(src reflect.Value) interface{}
	Transform(identifier string, handler TransformHandler) Mapper
	Ignore(identifier string) Mapper
	Condition(identifier string, conditionHandler ConditionHandler) Mapper
	Nested(identifier, target string, mapper Mapper) Mapper
}
```
* `Template` return the template of the mapper
* `Mapping` mapping object-object(the parameters need is a pointer), slice-slice, array-array
* `MakeDestination` like the `Mapping` method, but the destination object will be created inside this method
* `Transform` is the method define a custom transform for dst fields
* `Ignore` ignore mapping a field
* `Condition` execute `FromField` with default identifier if condition is true
* `Nested` execute `Mapping` from target field to identifier field via the parameter `mapper`

#### NameMapper
```Go
func NewNameMapper(template Template, profile map[string]string) (*NameMapper, error) 
func (nm *NameMapper) SetProfile(profile map[string]string) error 
```
NameMapper is the implement of the Mapper interface.\
Make a mapper via the `profile`, and the public fields' names.\
The `profile` is a `map[string]string` with keys are the names of the destination fields and values are the names of source fields are mapped together.\
The parameter `identifier` is the public name of the field.\
Default, the fields of the destination object will be mapping with fields of the source object via values defined in the `profile`.\
It `panic` if the type of the source field different type of the destination field.

## Template
```Go
func NewTemplate(src, dst interface{}) Template
func (t Template) ProfileSameTag() map[string]string
func (t Template) ProfileSameName() map[string]string 
```
* `NewTemplate` the method creates a Template. The `src` and `dst` parameters are template instances of structs; It can be an empty struct.
* `ProfileSameName` make a profile with the fields had the same name will be mapped to each other.
* `ProfileSameTag` make a profile with the fields had the same tag will be mapped to each other.

## TransformHandler
```Go
type TransformHandler func(SrcMap FieldMap) interface{}
````
TransformHandler is a callback with input is `FieldMap` of the source object,
and returns a field of the destination object.
```Go
func FromField(srcFieldName string) TransformHandler
````
This is the default handler. Will be returned to the target value.


## FieldMap
```Go
func (fm FieldMap) Mapping(SrcMap FieldMap, transformers map[string]TransformHandler) 
func (fm FieldMap) Field(identifier string) reflect.Value
func (f FieldMap)  String(fieldName string) string 
func (f FieldMap)  Int64(identifier string) int64
func (f FieldMap)  Int(fieldName string) int
func (f FieldMap)  Float(fieldName string) float64
func (f FieldMap)  Complex(identifier string) complex128
func (f FieldMap)  Bytes(fieldName string) []byte
func (f FieldMap)  Bool(fieldName string) bool
func (f FieldMap)  Interface(fieldName string) interface{}
```
FieldMap is the map of `reflect.Value` with a key is the field name.
You can call the `Field(fieldName string)` method or other methods to get the value of the field corresponding.\
You can't call the `Set` method to set the value for a variable
you get from FieldMap, it not references to the source object,
any attempt to change it will cause panic;
This is to make sure the source object will not be changed after mapping.\
The `Mapping` method will be mapping values from SrcMap to DstMap via the transform handlers corresponding.

## Util
```Go
func IGNORE() interface{} 
```
In the Transform methods,  return the result of the IGNORE() method to ignore mapping the current field.
## Example
```Go
package main

import (
	"automapper"
	"fmt"
)

type PersonUser struct {
	Username string `automapper:"user"`
}
type EmployeeUser struct {
	User string `automapper:"user"`
}

type Person struct {
	FirstName string       `automapper:"firstName"`
	LastName  string       `automapper:"lastName"`
	Age       int          `automapper:"age"`
	User      []PersonUser `automapper:"user"`
}

type Employee struct {
	FullName  string         `automapper:"fullName"`
	Age       int            `automapper:"age"`
	FirstName string         `automapper:"firstName"`
	LastName  string         `automapper:"lastName"`
	User      []EmployeeUser `automapper:"user"`
}

func main() {
	var err error
	//create test data
	LPerson := make([]Person, 2)
	LEmployee := make([]Employee, 2)
	OEmployee := Employee{}
	OPerson := Person{
		FirstName: "Hoang",
		LastName:  "Nguyen Hai",
		Age:       23,
		User: []PersonUser{
			{Username: "hoangnh01"},
			{Username: "hoangnh02"}},
	}
	LPerson[0] = Person{
		Age:       27,
		FirstName: "Thuoc",
		LastName:  "Nguyen Van",
		User: []PersonUser{
			{Username: "thuocnv01"},
			{Username: "thuocnv02"},
		},
	}
	LPerson[1] = Person{
		Age:       17,
		FirstName: "Quang Hung",
		LastName:  "Tran",
		User: []PersonUser{
			{Username: "hungqt01"},
			{Username: "hungqt02"},
		},
	}

	// new template and profile
	template := automapper.NewTemplate(Person{}, Employee{})
	userTemplate := automapper.NewTemplate(PersonUser{}, EmployeeUser{})

	// new mapper via tag
	mapper,err := automapper.NewNameMapper(template,template.ProfileSameTag())
	if err != nil{
		panic(err)
	}
	userMapper,err := automapper.NewNameMapper(userTemplate, userTemplate.ProfileSameName())
	if err != nil{
		panic(err)
	}
	mapper.Transform(
		"FullName", func(SrcMap automapper.FieldMap) interface{} {
			return SrcMap.String("FirstName") + " " + SrcMap.String("LastName")

		},
	).Ignore("LastName",
	).Condition("Age", func(SrcMap automapper.FieldMap) bool {
		return SrcMap.Int("Age") > 18

	}).Nested("User", "User", userMapper)

	// test via tag
	fmt.Println("______ViaTag_______")
	fmt.Println("______Struct_______")
	err = mapper.Mapping(&OPerson, &OEmployee)
	if err != nil {
		panic(err)
	}
	fmt.Println("Persons  :", OPerson)
	fmt.Println("Employees:", OEmployee)
	fmt.Println("______Slice_______")
	err = mapper.Mapping(LPerson, LEmployee)
	if err != nil {
		panic(err)
	}
	fmt.Println("ListPersons  :", LPerson)
	fmt.Println("ListEmployees:", LEmployee)

	// reset
	LEmployee = make([]Employee, 2)
	OEmployee = Employee{}

	// new mapper via profile
	mapper,err = automapper.NewNameMapper(template, template.ProfileSameName())
	if err != nil{
		panic(err)
	}
	userMapper,err = automapper.NewNameMapper(userTemplate, userTemplate.ProfileSameTag())
	if err != nil{
		panic(err)
	}
	mapper.Transform(
		"FullName", func(SrcMap automapper.FieldMap) interface{} {
			return SrcMap.String("FirstName") + " " + SrcMap.String("LastName")

		},
	).Ignore("LastName",
	).Condition("Age", func(SrcMap automapper.FieldMap) bool {
		return SrcMap.Int("Age") > 18

	}).Nested("User", "User", userMapper)

	// test via profile
	fmt.Println("______ViaProfile_______")
	fmt.Println("______Struct_______")
	err = mapper.Mapping(&OPerson, &OEmployee)
	if err != nil {
		panic(err)
	}
	fmt.Println("Persons  :", OPerson)
	fmt.Println("Employees:", OEmployee)
	fmt.Println("______Slice_______")
	err = mapper.Mapping(LPerson, LEmployee)
	if err != nil {
		panic(err)
	}
	fmt.Println("ListPersons  :", LPerson)
	fmt.Println("ListEmployees:", LEmployee)
}

```
