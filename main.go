// basic api for learning purposes, based on the animal api

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting the application...")

	handleRequests()
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "animaldb"
	dbPass := "dud123"
	dbName := "animals"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/animals", getAllAnimals).Methods("GET")
	myRouter.HandleFunc("/animals", createAnimal).Methods("POST")
	myRouter.HandleFunc("/animals/{id}", getAnimal).Methods("GET")
	myRouter.HandleFunc("/animals/{id}", deleteAnimal).Methods("DELETE")
	myRouter.HandleFunc("/animals/{id}", updateAnimal).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

type Animal struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
}

var Animals []Animal

func getAllAnimals(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getAllAnimals")
	db := dbConn()
	fetchAllAnimals(db)
	json.NewEncoder(w).Encode(Animals)
}

func fetchAllAnimals(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM animals")
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var animal Animal
		err = rows.Scan(&animal.ID, &animal.Name, &animal.Birthdate)
		if err != nil {
			panic(err.Error())
		}

		Animals = append(Animals, animal)
	}

	defer db.Close()
}

func getAnimal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getAnimal")
	vars := mux.Vars(r)
	id := vars["id"]
	db := dbConn()
	animal := fetchAnimal(db, id)
	json.NewEncoder(w).Encode(animal)
}

func fetchAnimal(db *sql.DB, id string) Animal {
	var animal Animal
	err := db.QueryRow("SELECT * FROM animals WHERE id = ?", id).Scan(&animal.ID, &animal.Name, &animal.Birthdate)
	if err != nil {
		panic(err.Error())
	}

	return animal
}

func createAnimal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createAnimal")
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	// fmt.Fprintf(w, "%+v", string(reqBody))
	var animal Animal
	json.Unmarshal(reqBody, &animal)
	Animals = append(Animals, animal)

	json.NewEncoder(w).Encode(animal)
}

func deleteAnimal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteAnimal")
	vars := mux.Vars(r)
	id := vars["id"]
	for index, animal := range Animals {
		if animal.ID == id {
			Animals = append(Animals[:index], Animals[index+1:]...)
		}
	}
}

func updateAnimal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateAnimal")
	vars := mux.Vars(r)
	id := vars["id"]
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var updatedAnimal Animal
	json.Unmarshal(reqBody, &updatedAnimal)

	for index, animal := range Animals {
		if animal.ID == id {
			Animals[index] = updatedAnimal
		}
	}
}
