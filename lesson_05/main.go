package main

import (
	"fmt"
	"lesson_05/documentstore"
	"lesson_05/users"
)

//import "lesson_05/documentstore"

func main() {

	store := documentstore.NewStore()
	_, collection := store.CreateCollection("users", &documentstore.CollectionConfig{PrimaryKey: "ID"})

	userService := users.NewService(*collection)
	userService.CreateUser("1", "Alice")
	userService.CreateUser("15", "Alice222")

	fmt.Println(userService.ListUsers())
	userService.DeleteUser("1")

	fmt.Println(userService.ListUsers())

	// -------------------------------------------
	//type Example struct {
	//	Name    string
	//	Age     float64
	//	IsAdmin bool
	//	Tags    []string
	//}
	//
	//doc := &documentstore.Document{
	//	Fields: map[string]documentstore.DocumentField{
	//		"Name":    {Type: documentstore.DocumentFieldTypeString, Value: "Alice"},
	//		"Age":     {Type: documentstore.DocumentFieldTypeNumber, Value: 25.0},
	//		"IsAdmin": {Type: documentstore.DocumentFieldTypeBool, Value: true},
	//		"Tags":    {Type: documentstore.DocumentFieldTypeArray, Value: []any{"go", "coding"}},
	//	},
	//}
	//
	//// Объект, в который будет происходить десериализация
	//var example Example
	//
	//// Вызов функции UnmarshalDocument
	//err := documentstore.UnmarshalDocument(doc, &example)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//
	//// Вывод десериализованных данных
	//fmt.Printf("Unmarshaled struct: %+v\n", example)

}
