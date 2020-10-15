# Automapper for Golang
## Introduction
The Automapper is small package enables for Go programs mapping value bettwen different struct

The import path for the package is github.com/gnoah1379/automapper\
To install it, run:
```
go get github.com/gnoah1379/automapper
```
## Using
#### Builder
Before mapping, you need add tag `automapper:"Field Name"` for all field can be mapping like that:
```Go
type Person struct {
	FirstName string `automapper:"firstName"`
	LastName  string `automapper:"lastName"`
	Age       int    `automapper:"age"`
}
```
And create a mapper for 2 this structure, using function: `NewBuilder(srcTemplate interface{}, dstTemplate interface{})`\
You need provide 2 parameters, it is 2 instance of your mapping struct, and it can be empty.\
`NewBuilder` return Builder interface, have 2 method:
```Go
type Builder interface {
	Set(fieldName string, transform TransformHandler) Builder
	Build() AutoMapper
}
```
> `Set` for setup custom transform handler,
> receive field name of destination struct and a transform handler corresponding.
> Default, Automapper will search field with a name corresponding in the source struct and transfer it to destination field,
> if this field does not exist destination field will be received default value.\
> You can call this function  multiple times before call function Build.
>
> `Build` return a new mapper
#### TransformHandler
```Go
type TransformHandler func(SrcMap FieldMap) interface{}
````
TransformHandler is a callback with input is `FieldMap` of source struct,
and returns a field of destination struct
#### FieldMap
FieldMap is the map of `reflect.Value` with key is field name.\
You can call method `Field(fieldName string)` for a get value of field corresponding.\
Or call these methods for a get exact value of field:
```Go
func (f FieldMap) String(fieldName string) string 
func (f FieldMap) Int(fieldName string) int
func (f FieldMap) Float(fieldName string) float64
func (f FieldMap) Bytes(fieldName string) []byte
func (f FieldMap) Bool(fieldName string) bool
func (f FieldMap) Interface(fieldName string) interface{}
```
You can't `Set` value for variable you get from FieldMap, it not references to source object, any attempt to change it will cause panic.\
This is to make sure the source object will not be changed after mapping.
## Example
```Go
package main

import (
	"fmt"
	"github.com/gnoah1379/automapper"
)

type User struct {
	Username string `automapper:"username"`
}

type Person struct {
	FirstName string `automapper:"firstName"`
	LastName  string `automapper:"lastName"`
	Age       int    `automapper:"age"`
	User      User   `automapper:"user"`
}

type Employee struct {
	FullName  string `automapper:"fullName"`
	FirstName string `automapper:"firstName"`
	LastName  string `automapper:"lastName"`
	User      User   `automapper:"user"`
	Age       int    `automapper:"age"`
}

func main() {
	mapper := automapper.NewBuilder(Person{}, Employee{},
	).Set("fullName", func(SrcMap automapper.FieldMap) interface{} {
		return SrcMap.String("firstName") + " " + SrcMap.String("lastName")
	}).Set("user", func(SrcMap automapper.FieldMap) interface{} {
		return User{
			Username: SrcMap.Interface("user").(User).Username + "@etc.vn",
		}
	}).Set("age", automapper.Ignore).Build()
	p := &Person{Age: 22, FirstName: "Nguyen Hai", LastName: "Hoang", User: User{Username: "hoangnh"}}
	e := new(Employee)
	err := mapper.Mapping(p, e)
	if err != nil{
		panic(err.Error())
	}
	fmt.Println(e)
  // &{Nguyen Hai Hoang Nguyen Hai Hoang {hoangnh@etc.vn} 0}

}
```
