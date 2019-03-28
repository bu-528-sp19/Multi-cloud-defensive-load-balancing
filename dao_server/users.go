package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type User struct {
	ID       int `json:"ID"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

var users []User

func GetUser(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(req)
	id, _ := strconv.Atoi(params["id"])
	user := getUserById(id)

	if user.ID == id {
		json.NewEncoder(w).Encode(user)
		return
	}
	json.NewEncoder(w).Encode(&User{})
}

func GetUsers(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	users := getUsers()
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	user = createUser(user)
	json.NewEncoder(w).Encode(user)
}
