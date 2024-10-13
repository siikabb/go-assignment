// basic api for learning purposes, based on the animal api

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting the application...")
	Animals = []Animal{
		Animal{ID: "1", Name: "Risto", SpeciesID: 1, Birthdate: "2019-01-01"},
		Animal{ID: "2", Name: "Tuomo", SpeciesID: 2, Birthdate: "2019-02-02"},
		Animal{ID: "3", Name: "Jari", SpeciesID: 1, Birthdate: "2019-03-03"},
	}
	handleRequests()
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
	SpeciesID int    `json:"species_id"`
	Birthdate string `json:"birthdate"`
}

var Animals []Animal

func getAllAnimals(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getAllAnimals")
	json.NewEncoder(w).Encode(Animals)
}

func getAnimal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getAnimal")
	vars := mux.Vars(r)
	key := vars["id"]

	for _, animal := range Animals {
		if animal.ID == key {
			json.NewEncoder(w).Encode(animal)
		}
	}
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
