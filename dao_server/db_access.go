package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

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

func gcpLoginRead() (*sql.DB, error) {
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
		return db, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}
	return db, err
}

func awsLoginRead() (*sql.DB, error) {
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
		return dbAWS, err
	}

	errAWS = dbAWS.Ping()
	if errAWS != nil {
		return dbAWS, err
	}

	return dbAWS, err
}

func dbLoginread() *sql.DB {
	//loop here
	//Get GCP host
	db, err = gcpLoginRead()
	err2 = db.Ping()
	if err != nil || err2 != nil {
		//get AWS host
		awsdb, awserr = awsLoginRead()
		awserr2 = awsdb.Ping()
		if awserr != nil || awserr2 != nil {
			panic(awserr)
		}
		return awsdb
	}
	return db

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

	///aws
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

	return db, dbAWS
}
