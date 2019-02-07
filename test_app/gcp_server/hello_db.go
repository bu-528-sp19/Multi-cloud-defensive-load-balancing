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
	testUser := getUser("jtrinh", "ec528")
	fmt.Println(testUser.Username, testUser.Email)
	userCars := getCarsForUser(testUser.Id)

	for _, car := range userCars {
		fmt.Println(car.Model)
		reservation := getReservationForCar(car.Id)
		fmt.Println(reservation.StartTime, reservation.EndTime)
	}

	garages := getGarages()

	for _, garage := range garages {
		fmt.Println(garage.Name)
		reservations := getReservationsForGarage(garage.Id)
		for _, reservation := range reservations {
			fmt.Println(reservation.StartTime, reservation.EndTime)
		}
	}
}

func getReservationsForGarage(garageID int) []Reservation {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query(
		"SELECT * FROM reservations WHERE garage_id = $1",
		garageID)

	var reservations []Reservation

	for rows.Next() {
		var id int
		var start string
		var end string
		var carID int
		var garageID int

		err = rows.Scan(&id, &start, &end, &carID, &garageID)
		if err != nil {
			panic(err)
		}

		reservations = append(reservations,
			Reservation{Id: id, StartTime: start, EndTime: end, CarId: carID, GarageId: garageID})
	}
	return reservations

}

func getReservationForCar(carID int) (Reservation) {
	db := dbLogin()
	defer db.Close()

	rows, err :=  db.Query(
		"SELECT * FROM reservations WHERE car_id = $1",
		carID)

	for rows.Next() {
		var id int
		var start string
		var end string
		var carID int
		var garageID int

		err = rows.Scan(&id, &start, &end, &carID, &garageID)
		if err != nil {
			panic(err)
		}

		return Reservation{Id: id, StartTime: start, EndTime: end, CarId: carID, GarageId: garageID}
	}
	return Reservation{}
}

func getGarages() ([]Garage) {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM garages")
	var garages []Garage

	for rows.Next() {
		var name string
		var maxCars int
		var id int

		err = rows.Scan(&id, &name, &maxCars)
		if err != nil {
			panic(err)
		}
		garages = append(garages, Garage{Id: id, Name: name, MaxCars: maxCars})
	}
	return garages
}

func getCarsForUser(userID int) ([]Car) {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query(
		"SELECT * FROM cars WHERE user_id = $1",
		userID)

	var userCars []Car
	for rows.Next() {
		var id int
		var user_id int
		var model string

		err = rows.Scan(&id, &user_id, &model)
		if err != nil {
			panic(err)
		}

		userCars = append(userCars, Car{Id: id, Model: model})
	}
	return userCars
}

func getUser(username string, password string) (User) {
	db := dbLogin()
	defer db.Close()

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

		return User{Id: id, Username: db_username, Email: email}
	}

	return User{}
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
	var dbInfo DatabaseInfo
	dbCredentials, _ := ioutil.ReadFile("db_secrets_gcp.json")
	json.Unmarshal(dbCredentials, &dbInfo)
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbInfo.Host,
		dbInfo.Port,
		dbInfo.User,
		dbInfo.Password,
		dbInfo.Name)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

