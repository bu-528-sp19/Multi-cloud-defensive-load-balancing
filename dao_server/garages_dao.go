package main

import (
	"time"
	"fmt"
)

func createGarage(garageObj Garage) Garage {
	query := fmt.Sprintf(
		"INSERT INTO garages (name, max_cars) "+
			"VALUES('%s', %d) "+
			"RETURNING id",
		garageObj.Name,
		garageObj.MaxCars)

	s.Set(time.Now().String(), query)

	db := dbLogin()
	defer db.Close()
	row, err := db.Query(query)

	if err != nil {
		panic (err)
	}

	row.Next()
	var newID int
	scanErr := row.Scan(&newID)

	if scanErr != nil {
		panic(scanErr)
	}

	garageObj.ID = newID
	return garageObj
}

func getGarages() ([]Garage) {
	db := dbLogin()
	defer db.Close()

	rows, _ := db.Query("SELECT * FROM garages")
	var garages []Garage

	for rows.Next() {
		var name string
		var maxCars int
		var id int

		err := rows.Scan(&id, &name, &maxCars)
		if err != nil {
			panic(err)
		}
		garages = append(garages, Garage{ID: id, Name: name, MaxCars: maxCars})
	}
	return garages
}

func getGarage(garageID int) (Garage) {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query(
		"SELECT * FROM garages WHERE id = $1",
		garageID)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int
		var name string
		var maxCars int

		scanErr := rows.Scan(&id, &name, &maxCars)
		if scanErr != nil {
			panic(err)
		}

		return Garage{ID: id, Name: name, MaxCars: maxCars}
	}
	return Garage{}
}

func deleteGarage(garageID int) {
	db := dbLogin()
	defer db.Close()

	_, err := db.Query(
		"DELETE from garages WHERE garages.id = $1",
		garageID)

	if err != nil {
		panic(err)
	}
}
