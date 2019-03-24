package main

import (
	"time"
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
)

const CARS_ROUTE string = "cars/"

func createCar(carObj Car) Car {
	if !s.IsLeader() {
		leaderIP := s.GetLeaderAddress() + ":8888"
		url := leaderIP + CARS_ROUTE

		jsonStr, _ := json.Marshal(carObj)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		_ = json.NewDecoder(resp.Body).Decode(&carObj)
		return carObj
	}

	query := fmt.Sprintf(
		"INSERT INTO cars (user_id, model) "+
			"VALUES (%d, '%s') RETURNING id;",
		carObj.UserID,
		carObj.Model)

	s.Set(time.Now().String(), query)

	db := dbLogin()
	defer  db.Close()
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

	carObj.ID = newID
	return carObj
}

func getCarsForUser(userID int) ([]Car) {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query(
		"SELECT * FROM cars WHERE user_id = $1",
		userID)

	var userCars []Car
	for rows.Next() {
		var id int
		var user_id int
		var model string

		err = rows.Scan(&id, &user_id, &model)
		if err != nil {
			panic(err)
		}

		userCars = append(userCars, Car{ID: id, Model: model, UserID: user_id})
	}
	return userCars
}

func getCar(carID int) (Car) {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query(
		"SELECT * FROM cars WHERE id = $1",
		carID)

	for rows.Next() {
		var id int
		var user_id int
		var model string

		err = rows.Scan(&id, &user_id, &model)
		if err != nil {
			panic(err)
		}
		return Car{ID: id, UserID: user_id, Model: model}
	}
	return Car{}
}

func getCars() ([]Car) {
	db := dbLogin()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM cars")

	var allCars []Car
	for rows.Next() {
		var id int
		var user_id int
		var model string

		err = rows.Scan(&id, &user_id, &model)
		if err != nil {
			panic(err)
		}
		allCars = append(allCars, Car{ID: id, UserID: user_id, Model: model})
	}
	return allCars
}

func deleteCar(carID int) {
	db := dbLogin()
	defer db.Close()

	_, err := db.Query(
		"DELETE FROM cars WHERE cars.id = $1",
		carID)

	if err != nil {
		panic (err)
	}
}
