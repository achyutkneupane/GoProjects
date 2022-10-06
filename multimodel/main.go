package main

import (
	"achyut.com.np/go/packs/multimodel/models"
)

func main() {
	println("Hello, world!")
	type User models.User
	type Todo models.Todo

	// set User
	var user User
	user.Id = 1
	user.Username = "achyut"
	user.Password = "achyut"
	user.Email = "achyutkneupane@gmail.com"
	user.FirstName = "Achyut"
	user.LastName = "Neupane"

	// set Todo
	var todo Todo
	todo.Id = 1
	todo.Title = "Learn Go"
	todo.Description = "Learn Go from scratch"
	
	println(user.Id)
	println(user.Username)
}

