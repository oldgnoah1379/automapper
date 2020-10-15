# Automapper for Golang
## Introduction
The Automapper is small package enables for Go programs mapping value bettwen different struct

The import path for the package is github.com/gnoah1379/automapper\
To install it, run:
```
go get github.com/gnoah1379/automapper
```
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
	err := mapper.Transform(p, e)
	if err != nil{
		panic(err.Error())
	}
	fmt.Println(e)
  // &{Nguyen Hai Hoang Nguyen Hai Hoang {hoangnh@etc.vn} 0}

}
```
