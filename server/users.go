package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID       string `json:"ID"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

var users []User

func GetUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range users {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

func GetUsers(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, req *http.Request) {
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	//@TODO make DB auto increment here
	user.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe
	users = append(users, user)
	json.NewEncoder(w).Encode(users)
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)
			user.ID = params["id"]
			users = append(users, user)
			json.NewEncoder(w).Encode(user)
			return
		}
	}
}
