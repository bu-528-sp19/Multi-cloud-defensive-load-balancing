package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Car struct {
	ID     string `json:"ID"`
	UserID string `json:"UserID"`
	Model  string `json:"Model"`
}

var cars []Car

func GetCar(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range cars {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Car{})
}

func GetCars(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(cars)
}

func CreateCar(w http.ResponseWriter, req *http.Request) {
	var car Car
	_ = json.NewDecoder(req.Body).Decode(&car)
	//@TODO make DB auto increment here
	car.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe
	cars = append(cars, car)
	json.NewEncoder(w).Encode(cars)
}

func DeleteCar(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range cars {
		if item.ID == params["id"] {
			cars = append(cars[:index], cars[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(cars)
}

func UpdateCar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range cars {
		if item.ID == params["id"] {
			cars = append(cars[:index], cars[index+1:]...)
			var car Car
			_ = json.NewDecoder(r.Body).Decode(&car)
			car.ID = params["id"]
			cars = append(cars, car)
			json.NewEncoder(w).Encode(car)
			return
		}
	}
}
