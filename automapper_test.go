package automapper

import "testing"

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

func TestAutoMapper_Transform(t *testing.T) {
	firstName := "Nguyen Hai"
	lastName := "Hoang"
	username := "hoangnh"
	company := "@etc.vn"
	mapper := NewBuilder(Person{}, Employee{}).Set("fullName", func(SrcMap FieldMap) interface{} {
		return SrcMap.String("firstName") + " " + SrcMap.String("lastName")
	}).Set("user", func(SrcMap FieldMap) interface{} {
		return User{
			Username: SrcMap.Interface("user").(User).Username + "@etc.vn",
		}
	}).Set("age", Ignore).Build()
	p := &Person{Age: 22, FirstName: firstName, LastName: lastName, User: User{Username: username}}
	e := new(Employee)
	err := mapper.Transform(p, e)
	if err != nil {
		t.Error(err.Error())
	}
	if e.FullName != firstName+" "+lastName {
		t.Error("FullName failed")
	}
	if e.User.Username != p.User.Username+company {
		t.Error("username failed")
	}
	if e.Age != 0 {
		t.Error("age failed")
	}
}
