// basic api for learning purposes, based on the animal api

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
