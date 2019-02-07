package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Garage struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	MaxCars string `json:"MaxCars"`
}

var garages []Garage

func GetGarage(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range garages {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Garage{})
}

func GetGarages(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(garages)
}

func CreateGarage(w http.ResponseWriter, req *http.Request) {
	var garage Garage
	_ = json.NewDecoder(req.Body).Decode(&garage)
	//@TODO make DB auto increment here
	garage.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe
	garages = append(garages, garage)
	json.NewEncoder(w).Encode(garages)
}

func DeleteGarage(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range garages {
		if item.ID == params["id"] {
			garages = append(garages[:index], garages[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(garages)
}

func UpdateGarage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range garages {
		if item.ID == params["id"] {
			garages = append(garages[:index], garages[index+1:]...)
			var garage Garage
			_ = json.NewDecoder(r.Body).Decode(&garage)
			garage.ID = params["id"]
			garages = append(garages, garage)
			json.NewEncoder(w).Encode(garage)
			return
		}
	}
}
