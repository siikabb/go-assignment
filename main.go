package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Starting the application...")

	handleRequests()
}

type Animal struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
}

type AnimalWithoutID struct {
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
}

var Animals []Animal
