package main

import (
	"fmt"
	"lesson_05/documentstore"
	"lesson_05/users"
)

func main() {

	store := documentstore.NewStore()
	collection, _ := store.CreateCollection("users", &documentstore.CollectionConfig{PrimaryKey: "ID"})
	userService := users.NewService(*collection)

	userId := "1"
	userName := "Alice"

	_, err := userService.CreateUser(userId, userName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	userList, err := userService.ListUsers()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(userList)

	user, err := userService.GetUser("1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(user)

	user, err = userService.GetUser("2")
	if err != nil {
		fmt.Println(err.Error())
	}

	err = userService.DeleteUser("2")
	if err != nil {
		fmt.Println(err.Error())
	}

	err = userService.DeleteUser("1")
	if err != nil {
		fmt.Println(err.Error())
	}

	userList, err = userService.ListUsers()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(userList)
}
