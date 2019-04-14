package main

import (
	"time"
	"fmt"
	"strconv"
	"net/http"
	"bytes"
	"encoding/json"
	"strings"
)

const RESERVATIONS_ROUTE string = "reservations/"

func createReservation(reservationObj Reservation) Reservation {

	//local server
	/*isLeader := true
	if isLeader == false { //!s.IsLeader() {
		leaderIP := "localhost:8888"*/

	//cloud server
	if !s.IsLeader() {
		leaderIP := "http://" + strings.Split(s.GetLeaderAddress(), ":")[0] + ":8888/"

		url := leaderIP + RESERVATIONS_ROUTE
		//fmt.Println(url)
		jsonStr, _ := json.Marshal(reservationObj)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		_ = json.NewDecoder(resp.Body).Decode(&reservationObj)
		return reservationObj
	}

  num_dbs := 0
	dbs := dbLogin()
	for _, db := range dbs {
		num_dbs = num_dbs + 1
		defer  db.Close()
	}

  time_layout := "2006-01-02T15:04:05.000Z"
	start, _ := time.Parse(time_layout, reservationObj.StartTime)
	end, _ := time.Parse(time_layout, reservationObj.EndTime)

	query := fmt.Sprintf(
		"INSERT INTO reservations (start_time, end_time, car_id, garage_id) "+
			"VALUES('%s', '%s', %d, %d) "+
			"RETURNING id",
		start,
		end,
		reservationObj.CarID,
		reservationObj.GarageID)

  cur_time := strconv.FormatInt(time.Now().Unix(), 10)
	s.Set(cur_time, query)

	if num_dbs == 0 {
		return reservationObj
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

	reservationObj.ID = newID
	return reservationObj
}

func getReservationsForUser(userID int) []Reservation {
  db, err := dbLoginread()
	defer db.Close()

	var reservations []Reservation
	//return no dbs working
	if err != nil {
		return reservations
	}

	rows, err := db.Query(
		"SELECT r.id, r.start_time, r.end_time, r.car_id, r.garage_id "+
		"FROM reservations r "+
		"INNER JOIN cars on r.car_id = cars.id "+
		"INNER JOIN users on users.id = cars.user_id "+
		"WHERE users.id = $1",
		userID)

	for rows.Next() {
		var id int
		var start string
		var end string
		var carID int
		var garageID int

		err = rows.Scan(&id, &start, &end, &carID, &garageID)
		if err != nil {
			panic(err)
		}

		reservations = append(reservations,
			Reservation{ID: id, StartTime: start, EndTime: end, CarID: carID, GarageID: garageID})
	}
	return reservations

}

func getReservationsForGarage(garageID int) []Reservation {
  db, err := dbLoginread()
	defer db.Close()

	var reservations []Reservation
	//return no dbs working
	if err != nil {
		return reservations
	}

	rows, err := db.Query(
		"SELECT * FROM reservations WHERE garage_id = $1",
		garageID)

	for rows.Next() {
		var id int
		var start string
		var end string
		var carID int
		var garageID int

		err = rows.Scan(&id, &start, &end, &carID, &garageID)
		if err != nil {
			panic(err)
		}

		reservations = append(reservations,
			Reservation{ID: id, StartTime: start, EndTime: end, CarID: carID, GarageID: garageID})
	}
	return reservations

}

func getReservationForCar(carID int) (Reservation) {
  db, err := dbLoginread()
	defer db.Close()

	//return no dbs working
	if err != nil {
		return Reservation{}
	}

	rows, err :=  db.Query(
		"SELECT * FROM reservations WHERE car_id = $1",
		carID)

	for rows.Next() {
		var id int
		var start string
		var end string
		var carID int
		var garageID int

		err = rows.Scan(&id, &start, &end, &carID, &garageID)
		if err != nil {
			panic(err)
		}

		return Reservation{ID: id, StartTime: start, EndTime: end, CarID: carID, GarageID: garageID}
	}
	return Reservation{}
}

func getReservations() ([]Reservation) {
  db, err := dbLoginread()
	defer db.Close()

	var reservations []Reservation
	//return no dbs working
	if err != nil {
		return reservations
	}

	rows, err := db.Query("SELECT * FROM reservations")

	for rows.Next() {
		var id int
		var start string
		var end string
		var carID int
		var garageID int

		err = rows.Scan(&id, &start, &end, &carID, &garageID)
		if err != nil {
			panic(err)
		}

		reservations = append(reservations,
			Reservation{ID: id, StartTime: start, EndTime: end, CarID: carID, GarageID: garageID})
	}

	return reservations
}

func getReservation(reservationID int) (Reservation) {
  db, err := dbLoginread()
	defer db.Close()

	//return no dbs working
	if err != nil {
		return Reservation{}
	}

	rows, err := db.Query(
		"SELECT * FROM reservations WHERE id = $1",
		reservationID)

	for rows.Next() {
		var id int
		var start string
		var end string
		var carID int
		var garageID int

		err = rows.Scan(&id, &start, &end, &carID, &garageID)
		if err != nil {
			panic(err)
		}

		return Reservation{ID: id, StartTime: start, EndTime: end, CarID: carID, GarageID: garageID}
	}
	return Reservation{}
}

/*func deleteReservation(reservationID int) {
	db,db2 := dbLogin()
	defer db.Close()
	defer db2.Close()

	_, err := db.Query(
		"DELETE FROM reservations WHERE reservations.id = $1",
		reservationID)

	_, err2 := db2.Query(
		"DELETE FROM reservations WHERE reservations.id = $1",
		reservationID)

	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}
}*/
