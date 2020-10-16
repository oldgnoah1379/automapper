# Automapper for Golang
## Introduction
The Automapper is the package for Go programs mapping values in different structs

The import path for the package is `github.com/gnoah1379/automapper`\
To install it, run:
```
go get github.com/gnoah1379/automapper
```
## Using
#### Builder
First, you need to add tag `automapper:"Field Name"` for all fields can be mapping like that:
```Go
type Person struct {
	FirstName string `automapper:"firstName"`
	LastName  string `automapper:"lastName"`
	Age       int    `automapper:"age"`
}
```
And create a mapper for 2 this structure, using function: `NewBuilder(srcTemplate interface{}, dstTemplate interface{})`\
You need to provide 2 parameters, it is 2 instances of your mapping struct, and it can be empty.\
`NewBuilder` return Builder interface, have 2 methods:
```Go
type Builder interface {
	Set(fieldName string, transform TransformHandler) Builder
	Build() AutoMapper
}
```
> `Set` for setup custom transform handler receive 2 parameters,\
> fist param is FieldName of destination object\
> second param is a transform handler corresponding.
> Default, Automapper will search field with a name corresponding in the source object and transfer it to destination field,
> if this field does not exist, the destination field is default value.\
> You can call this function  multiple times before call the Build function.
>
> `Build` return a new mapper
#### TransformHandler
```Go
type TransformHandler func(SrcMap FieldMap) interface{}
````
TransformHandler is a callback with input is `FieldMap` of the source object,
and returns a field of the destination object.\
###### Handler
`automapper.Ignore` will be ignored mapping a field corresponding.\
 `automapper.Default` is default handler for all fields. 
###### Util
`automapper.Condition` first param is condition callback `func(SrcMap FieldMap) bool `\
if condition equal true, will execute TransformHandler is second param(if this param is nil, will execute Default),
if condition equal false, will execute Ignore handler

#### FieldMap
FieldMap is the map of `reflect.Value` with a key is the field name.
You can call the `Field(fieldName string)` method to get the value of the field corresponding.
Or call these methods to get the exact value of the field:
```Go
func (f FieldMap) String(fieldName string) string 
func (f FieldMap) Int(fieldName string) int
func (f FieldMap) Float(fieldName string) float64
func (f FieldMap) Bytes(fieldName string) []byte
func (f FieldMap) Bool(fieldName string) bool
func (f FieldMap) Interface(fieldName string) interface{}
```
You can't call the `Set` method for set the value for a variable
you get from FieldMap, it not references to the source object,
any attempt to change it will cause panic;
This is to make sure the source object will not be changed
after mapping.
#### AutoMapper
`Mapping(src interface{}, dst interface{}) error` 
for mapping object-to-object\
`ListMapping(src interface{}, dst interface{}) error` 
for mapping slice-to-slice
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
