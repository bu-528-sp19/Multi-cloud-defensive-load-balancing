package main

import (
	"encoding/json"
	"net/http"
)

type History struct {
	ID            string `json:"ID"`
	ReservationID string `json:"ReservationID"`
}

var history []History

func GetHistory(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(history)
}
