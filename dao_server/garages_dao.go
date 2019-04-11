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

const GARAGES_ROUTE string = "garages/"

func createGarage(garageObj Garage) Garage {

	//local server
	/*isLeader := true
	if isLeader == false { //!s.IsLeader() {
		leaderIP := "localhost:8888"*/

	//cloud server
	if !s.IsLeader() {
		leaderIP := "http://" + strings.Split(s.GetLeaderAddress(), ":")[0] + ":8888/"

		url := leaderIP + GARAGES_ROUTE
		//fmt.Println(url)
		jsonStr, _ := json.Marshal(garageObj)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		_ = json.NewDecoder(resp.Body).Decode(&garageObj)
		return garageObj
	}

  num_dbs := 0
	dbs := dbLogin()
	for _, db := range dbs {
		num_dbs = num_dbs + 1
		defer  db.Close()
	}

	query := fmt.Sprintf(
		"INSERT INTO garages (name, max_cars) "+
			"VALUES('%s', %d) "+
			"RETURNING id",
		garageObj.Name,
		garageObj.MaxCars)

  cur_time := strconv.FormatInt(time.Now().Unix(), 10)
	s.Set(cur_time, query)

  if num_dbs == 0 {
		return garageObj
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

	garageObj.ID = newID
	return garageObj
}

func getGarages() ([]Garage) {
  db, err := dbLoginread()
	defer db.Close()

	var garages []Garage
	//return no dbs working
	if err != nil {
		return garages
	}

	rows, _ := db.Query("SELECT * FROM garages")

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
  db, err := dbLoginread()
	defer db.Close()

	//return no dbs working
	if err != nil {
		return Garage{}
	}

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

/*
func deleteGarage(garageID int) {
	db,db2 := dbLogin()
	defer db.Close()
	defer db2.Close()

	_, err := db.Query(
		"DELETE from garages WHERE garages.id = $1",
		garageID)
	_, err2 := db2.Query(
		"DELETE from garages WHERE garages.id = $1",
		garageID)

	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}
}
*/
