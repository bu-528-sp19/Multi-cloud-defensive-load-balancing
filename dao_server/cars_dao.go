package main

import (
	"time"
	"fmt"
	"strconv"
	"net/http"
	"bytes"
	"encoding/json"
	//"strings"
)

const CARS_ROUTE string = "cars/"

func createCar(carObj Car) Car {

	isLeader := true

	if isLeader == false { //!s.IsLeader() {
		//		"http://" + strings.Split(s.GetLeaderAddress(), ":")[0] + ":8888/"
		leaderIP := "localhost:8888"//"http://" + strings.Split(s.GetLeaderAddress(), ":")[0] + ":8888/"
		url := leaderIP + CARS_ROUTE
		//fmt.Println(url)
		jsonStr, _ := json.Marshal(carObj)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		_ = json.NewDecoder(resp.Body).Decode(&carObj)
		return carObj
	}


	num_dbs := 0
	dbs := dbLogin()
	for _, db := range dbs {
		num_dbs = num_dbs + 1
		defer  db.Close()
	}

	query := fmt.Sprintf(
		"INSERT INTO cars (user_id, model) "+
		"VALUES (%d, '%s') RETURNING id;",
		carObj.UserID, carObj.Model)

	cur_time := strconv.FormatInt(time.Now().Unix(), 10)
	s.Set(cur_time, query)

	if num_dbs == 0 {
		return carObj
	}

	//execute queries to all dbs except the last one
	var i int
	for i = 0; i < num_dbs - 1; i++ {
		_, err := dbs[i].Query(query)
		if err != nil {
			panic (err)
		}
	}

	//only get data from last query
	row, err := dbs[len(dbs) - 1].Query(query)

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
	db, err := dbLoginread()
	defer db.Close()

	var userCars []Car
	//return no dbs working
	if err != nil {
		return userCars
	}

	rows, err := db.Query(
		"SELECT * FROM cars WHERE user_id = $1",
		userID)

	//var userCars []Car
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
	db, err := dbLoginread()
	defer db.Close()

	//return no dbs working
	if err != nil {
		return Car{}
	}

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
	db, err := dbLoginread()
	defer db.Close()

	var allCars []Car
	if err != nil{
		return allCars
	}

	rows, err := db.Query("SELECT * FROM cars")

	//var allCars []Car
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
