
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

const USERS_ROUTE string = "users/"

func createUser(userObj User) User {

	//isLeader := true

	if !s.IsLeader() {
		//		"http://" + strings.Split(s.GetLeaderAddress(), ":")[0] + ":8888/"
		leaderIP := "http://" + strings.Split(s.GetLeaderAddress(), ":")[0] + ":8888/"
		url := leaderIP + USERS_ROUTE
		//fmt.Println(url)
		jsonStr, _ := json.Marshal(userObj)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()
		_ = json.NewDecoder(resp.Body).Decode(&userObj)
		return userObj
	}

	num_dbs := 0
	dbs := dbLogin()
	for _, db := range dbs {
		num_dbs = num_dbs + 1
		defer  db.Close()
	}

	query := fmt.Sprintf(
		"INSERT INTO users (username, password, email) "+
			"VALUES ('%s', '%s', '%s') "+
			"RETURNING id",
		userObj.Username,
		userObj.Password,
		userObj.Email)

	cur_time := strconv.FormatInt(time.Now().Unix(), 10)
	s.Set(cur_time, query)

	if num_dbs == 0 {
		return userObj
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

	userObj.ID = newID
	return userObj
}

func getUserById(userID int) (User) {
	db, err := dbLoginread()
	defer db.Close()

	//return no dbs working
	if err != nil {
		return User{}
	}

	rows, err := db.Query(
		"SELECT * FROM users WHERE id = $1",
		userID)

	for rows.Next() {
		var id int
		var db_username string
		var password string
		var email string

		err = rows.Scan(&id, &db_username, &password, &email)
		if err != nil {
			panic(err)
		}

		return User{ID: id, Username: db_username, Email: email}
	}

	return User{}
}

func getUsers() ([]User) {
	db, err := dbLoginread()
	defer db.Close()

	var users []User
	//return no dbs working
	if err != nil {
		return users
	}

	rows, err := db.Query("SELECT * FROM users")

	for rows.Next() {
		var id int
		var db_username string
		var password string
		var email string

		err = rows.Scan(&id, &db_username, &password, &email)
		if err != nil {
			panic (err)
		}
		users = append(users,
			User{ID: id, Username: db_username, Email: email})
	}
	return users
}

func getUser(username string, password string) (User) {
	db, err := dbLoginread()
	defer db.Close()

	//return no dbs working
	if err != nil {
		return User{}
	}

	rows, err := db.Query(
		"SELECT * FROM users WHERE username = $1 AND password = $2",
		username,
		password)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int
		var db_username string
		var password string
		var email string

		err = rows.Scan(&id, &db_username, &password, &email)
		if err != nil {
			panic(err)
		}

		return User{ID: id, Username: db_username, Email: email}
	}

	return User{}
}

/*
func deleteUser(userID int) {
	db,db2 := dbLogin()
	defer db.Close()
	defer db2.Close()

	_, err := db.Query(
		"DELETE FROM users where users.id = $1",
		userID)
	_, err2 := db2.Query(
		"DELETE FROM users where users.id = $1",
		userID)
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}
}
*/
