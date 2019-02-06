package main

import (
  "database/sql"
  "fmt"
  "encoding/json"
  "io/ioutil"
  _ "github.com/lib/pq"
)

type DatabaseInfo struct {
	Host string
	User string
	Password string
	Name string
	Port int
}

func main() {

  var dbInfo DatabaseInfo
  dbCredentials, _ := ioutil.ReadFile("db_secrets_gcp.json")
  fmt.Println(dbCredentials)
  json.Unmarshal(dbCredentials, &dbInfo)
  fmt.Println(dbInfo.Host, dbInfo.Port, dbInfo.User, dbInfo.Password, dbInfo.Name)
psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    dbInfo.Host, dbInfo.Port, dbInfo.User, dbInfo.Password, dbInfo.Name)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }

  fmt.Println("Successfully connected!")

  rows, err := db.Query("SELECT * FROM users")
  if err != nil {
	fmt.Println(err)
  }

  for rows.Next() {
    var id int
	var username string
    var password string
	var email string

	err = rows.Scan(&id, &username, &password, &email)
	fmt.Println(id, username, password, email)
  }
}
