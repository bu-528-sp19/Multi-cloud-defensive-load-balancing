package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Reservation struct {
	ID        string `json:"ID"`
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
	CarID     string `json:"CarID"`
	GarageID  string `json:"GarageID"`
}

var reservations []Reservation

func GetReservation(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range reservations {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Reservation{})
}

func GetReservations(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(reservations)
}

func CreateReservation(w http.ResponseWriter, req *http.Request) {
	var reservation Reservation
	_ = json.NewDecoder(req.Body).Decode(&reservation)
	//@TODO make DB auto increment here
	reservation.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe
	reservations = append(reservations, reservation)
	json.NewEncoder(w).Encode(reservations)
}

func DeleteReservation(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range reservations {
		if item.ID == params["id"] {
			reservations = append(reservations[:index], reservations[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(reservations)
}

func UpdateReservation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range reservations {
		if item.ID == params["id"] {
			reservations = append(reservations[:index], reservations[index+1:]...)
			var reservation Reservation
			_ = json.NewDecoder(r.Body).Decode(&reservation)
			reservation.ID = params["id"]
			reservations = append(reservations, reservation)
			json.NewEncoder(w).Encode(reservation)
			return
		}
	}
}
