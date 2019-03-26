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

	cur_time := time.Now().String()
	s.Set(cur_time, query)

	//db := dbLogin() // dbLogin now returns both
	db,db2 := dbLogin()
	defer  db.Close()
	defer  db2.Close()

	//////////////////////////////////////////////////
	//Database Forwarding

	conn_err := db.Ping()
	if conn_err != nil {
		aws_status, _ := s.Get("AWS_DOWN")
		if aws_status == "0" {
			s.Set("AWS_DOWN", cur_time)
		}
	} else{
			down_time, _ := s.Get("AWS_DOWN")
			if down_time != "0" {
				s.Set("AWS_DOWN", "0")
				db_recover(db, down_time)
			}
	}
	//////////////////////////////////////////////////

	row, err := db.Query(query)
	_, err2 := db2.Query(query)

	if err != nil {
		panic (err)
	}
	if err2 != nil {
		panic (err2)
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
	db := dbLoginread() //now dbLoginread instead of dbLogin
	defer db.Close()

	//////////////////////////////////////////////////
	//Database Forwarding

	if len(s.GetAll()) == 0 {
		s.Set("AWS_DOWN", "0")
	}

	isLeader := "true"

	conn_err := db.Ping()
	if conn_err != nil {
		fmt.Println("i guess the db is down\n")
		aws_status, _ := s.Get("AWS_DOWN")
		if aws_status == "0" {
			if isLeader != "true" {
				notify_leader("AWS_DOWN", time.Now().String())
			} else{
				s.Set("AWS_DOWN", time.Now().String())
			}
		}
		//read from other db
	} else{
			down_time, _ := s.Get("AWS_DOWN")
			if down_time != "0" {
				fmt.Println("should not be here unless db was down\n")
				if isLeader != "true" {
					notify_leader("AWS_DOWN", "0")
				} else{
					s.Set("AWS_DOWN", "0")
				}
				db_recover(db, down_time)
			}
	}
	//////////////////////////////////////////////////

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
	db := dbLoginread()
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
	db := dbLoginread()
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
	db,db2 := dbLogin()
	defer db.Close()
	defer db2.Close()

	_, err := db.Query(
		"DELETE FROM cars WHERE cars.id = $1",
		carID)
	_, err2 := db2.Query(
		"DELETE FROM cars WHERE cars.id = $1",
		carID)

	if err != nil {
		panic (err)
	}
	if err2 != nil {
		panic (err2)
	}
}
