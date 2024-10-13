package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

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

func fetchAnimal(db *sql.DB, id string) Animal {
	var animal Animal
	err := db.QueryRow("SELECT * FROM animals WHERE id = ?", id).Scan(&animal.ID, &animal.Name, &animal.Birthdate)
	if err != nil {
		panic(err.Error())
	}

	return animal
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
