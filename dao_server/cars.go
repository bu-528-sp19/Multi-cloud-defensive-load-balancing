package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Car struct {
	ID     int `json:"ID"`
	UserID int `json:"UserID"`
	Model  string `json:"Model"`
}

func GetCar(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	id, _ := strconv.Atoi(params["id"])

	car := getCar(id)
	if car.ID == id {
		json.NewEncoder(w).Encode(car)
		return
	}
	json.NewEncoder(w).Encode(&Car{})
}

func GetCars(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	cars := getCars()
	json.NewEncoder(w).Encode(cars)
}

func GetCarsByUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	id, _ := strconv.Atoi(params["id"])

	cars := getCarsForUser(id)
	json.NewEncoder(w).Encode(cars)
}

func CreateCar(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	var car Car
	_ = json.NewDecoder(req.Body).Decode(&car)
	car = createCar(car)
	json.NewEncoder(w).Encode(car)
}
