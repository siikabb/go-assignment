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

type AnimalWithoutID struct {
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

	Animals = []Animal{}

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

	var animal AnimalWithoutID
	json.Unmarshal(reqBody, &animal)
	db := dbConn()
	insertAnimal(db, animal)
	json.NewEncoder(w).Encode(animal)
}

func insertAnimal(db *sql.DB, animal AnimalWithoutID) {
	insert, err := db.Prepare("INSERT INTO animals (name, birthdate) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = insert.Exec(animal.Name, animal.Birthdate)
	if err != nil {
		panic(err.Error())
	}
}

func deleteAnimal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteAnimal")
	vars := mux.Vars(r)
	id := vars["id"]
	db := dbConn()
	removeAnimal(db, id)
}

func removeAnimal(db *sql.DB, id string) {
	delete, err := db.Prepare("DELETE FROM animals WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = delete.Exec(id)
	if err != nil {
		panic(err.Error())
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

	var animal Animal
	json.Unmarshal(reqBody, &animal)
	animal.ID = id
	db := dbConn()
	editAnimal(db, animal)
	json.NewEncoder(w).Encode(animal)

}

func editAnimal(db *sql.DB, animal Animal) {
	update, err := db.Prepare("UPDATE animals SET name = ?, birthdate = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = update.Exec(animal.Name, animal.Birthdate, animal.ID)
	if err != nil {
		panic(err.Error())
	}
}
