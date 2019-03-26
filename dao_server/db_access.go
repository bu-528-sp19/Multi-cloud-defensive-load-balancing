package main

import (
  "database/sql"
  "fmt"
  "os"
  "strconv"
  "math/rand"
  "time"
  _ "github.com/lib/pq"
)

type DatabaseInfo struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     int
}

//deleted test select function

func dbLoginread() (*sql.DB) {
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
	}else if (errAWS!=nil)&&(err==nil){
		return db
	}else {
		rand.Seed(time.Now().UnixNano())
		dice :=rand.Intn(20)
		if (dice%2 == 0){
			return db
		}else{
			return dbAWS
		}

	}

}

func dbLogin() (*sql.DB, *sql.DB) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"),
		port,
		os.Getenv("USER"),
		os.Getenv("PASSWORD"),
		os.Getenv("NAME"))
	//db, err := sql.Open("postgres", psqlInfo)//err declared and not used
	db, _ := sql.Open("postgres", psqlInfo)
	port, _ = strconv.Atoi(os.Getenv("PORT"))
	psqlInfoAWS := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HOSTAWS"),
		port,
		os.Getenv("USER"),
		os.Getenv("PASSWORDAWS"),
		os.Getenv("NAME"))
	//dbAWS, errAWS := sql.Open("postgres", psqlInfoAWS)//errAWs declared and not used
	dbAWS, _ := sql.Open("postgres", psqlInfoAWS)

	return db,dbAWS
}
