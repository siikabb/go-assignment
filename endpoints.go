package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

func getAllAnimals(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getAllAnimals")
	db := dbConn()
	fetchAllAnimals(db)
	json.NewEncoder(w).Encode(Animals)
}

func getAnimal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getAnimal")
	vars := mux.Vars(r)
	id := vars["id"]
	db := dbConn()
	animal := fetchAnimal(db, id)
	json.NewEncoder(w).Encode(animal)
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

func deleteAnimal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteAnimal")
	vars := mux.Vars(r)
	id := vars["id"]
	db := dbConn()
	removeAnimal(db, id)
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
