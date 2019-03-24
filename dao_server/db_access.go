package main

import (
  "database/sql"
  "fmt"
  "os"
  "strconv"
  "math/rand"
  _ "github.com/lib/pq"
)

type DatabaseInfo struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     int
}

func testSelect() {
	db := dbLogin()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int
		var username string
		var password string
		var email string

		err = rows.Scan(&id, &username, &password, &email)
		if err != nil {
			panic(err)
		}

		fmt.Println(id, username, password, email)
	}
}

func dbLogin() (*sql.DB) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"),
		port,
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("NAME"))
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}



//GCP

	port, _ = strconv.Atoi(os.Getenv("PORT"))
	psqlInfoAWS := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOSTAWS"),
		port,
		os.Getenv("USER"),
		os.Getenv("PASSWORDAWS"),
		os.Getenv("NAME"))
	dbAWS, errAWS := sql.Open("postgres", psqlInfoAWS)

	if errAWS != nil {
		panic(errAWS)
	}

	errAWS = dbAWS.Ping()
	if errAWS != nil {
		panic(errAWS)
	}


	if (errAWS == nil) && (err!=nil){
		return dbAWS
	}
	else if(errAWS!=nil)&&(err==nil){
		return db
	}
	else if(errAWS==nil)&&(err==nil){
		rand.Seed(time.Now().UnixNano())
		dice=rand.Intn(20)
		if (dice%2 == 0){
			return db
		}
		else{
			return dbAWS
		}

	}


	
}
