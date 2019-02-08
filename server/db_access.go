package main

import (
  "database/sql"
  "fmt"
  "encoding/json"
  "io/ioutil"
  "time"
  _ "github.com/lib/pq"
)

type DatabaseInfo struct {
	Host string
	User string
	Password string
	Name string
	Port int
}

/*func main() {
	testUser := getUser("jtrinh", "ec528")
	fmt.Println(testUser.Username, testUser.Email)
	userCars := getCarsForUser(testUser.ID)

	for _, car := range userCars {
		fmt.Println(car.Model)
		reservation := getReservationForCar(car.ID)
		fmt.Println(reservation.StartTime, reservation.EndTime)
	}

	garages := getGarages()

	for _, garage := range garages {
		fmt.Println(garage.Name)
		reservations := getReservationsForGarage(garage.ID)
		for _, reservation := range reservations {
			fmt.Println(reservation.StartTime, reservation.EndTime)
		}
	}

	reservationsForUser := getReservationsForUser(testUser.ID)
	for _, reservation := range reservationsForUser {
		fmt.Println(reservation.StartTime, reservation.EndTime)
	}

	newUser := User{Username: "TestUser", Password: "lolhashed", Email: "pls@bu.edu"}
	newUser = createUser(newUser)

	fmt.Println("User ID: ", newUser.ID)

	newCar := Car{Model: "R8", UserID: newUser.ID}
	newCar = createCar(newCar)

	newGarage := Garage{Name: "Decent", MaxCars: 7}
	newGarage = createGarage(newGarage)

	newReservation := Reservation {StartTime: time.Now().String(), EndTime: time.Now().String(), CarID: newCar.ID,GarageID: newGarage.ID}
	newReservation = createReservation(newReservation)

	deleteReservation(newReservation.ID)
	deleteGarage(newGarage.ID)
	deleteCar(newCar.ID)
	deleteUser(newUser.ID)
}*/

func deleteReservation(reservationID int) {
	db := dbLogin()
	defer db.Close()

	_, err := db.Query(
		"DELETE FROM reservations WHERE reservations.id = $1",
		reservationID)

	if err != nil {
		panic(err)
	}
}

func createReservation(reservationObj Reservation) Reservation {
	db := dbLogin()
	defer db.Close()

	time_layout := "2006-01-02T15:04:05.000Z"
	start, _ := time.Parse(time_layout, reservationObj.StartTime)
	end, _ := time.Parse(time_layout, reservationObj.EndTime)
	row, err := db.Query(
		"INSERT INTO reservations (start_time, end_time, car_id, garage_id) "+
		"VALUES($1, $2, $3, $4) "+
		"RETURNING id",
		start,
		end,
		reservationObj.CarID,
		reservationObj.GarageID)

	if err != nil {
		panic(err)
	}

	row.Next()
	var newID int
	scanErr := row.Scan(&newID)

	if scanErr != nil {
		panic(scanErr)
	}

	reservationObj.ID = newID
	return reservationObj
}

func deleteGarage(garageID int) {
	db := dbLogin()
	defer db.Close()

	_, err := db.Query(
		"DELETE from garages WHERE garages.id = $1",
		garageID)

	if err != nil {
		panic(err)
	}
}

func createGarage(garageObj Garage) Garage {
	db := dbLogin()
	defer db.Close()

	row, err := db.Query(
		"INSERT INTO garages (name, max_cars) "+
		"VALUES($1, $2) "+
		"RETURNING id",
		garageObj.Name,
		garageObj.MaxCars)

	if err != nil {
		panic (err)
	}

	row.Next()
	var newID int
	scanErr := row.Scan(&newID)

	if scanErr != nil {
		panic(scanErr)
	}

	garageObj.ID = newID
	return garageObj
}

func deleteCar(carID int) {
	db := dbLogin()
	defer db.Close()

	_, err := db.Query(
		"DELETE FROM cars WHERE cars.id = $1",
		carID)

	if err != nil {
		panic (err)
	}
}

func createCar(carObj Car) Car {
	db := dbLogin()
	defer  db.Close()

	row, err := db.Query(
		"INSERT INTO cars (user_id, model) "+
		"VALUES ($1, $2) "+
		"RETURNING id",
		carObj.UserID,
		carObj.Model)

	if err != nil {
		panic (err)
	}

	row.Next()
	var newID int
	scanErr := row.Scan(&newID)

	if scanErr != nil {
		panic(scanErr)
	}

	carObj.ID = newID
	return carObj
}

func deleteUser(userID int) {
	db := dbLogin()
	defer db.Close()

	_, err := db.Query(
		"DELETE FROM users where users.id = $1",
		userID)

	if err != nil {
		panic(err)
	}
}

func createUser(userObj User) User {
	db := dbLogin()
	defer db.Close()

	row, err := db.Query(
		"INSERT INTO users (username, password, email) "+
		"VALUES ($1, $2, $3) "+
		"RETURNING id",
		userObj.Username,
		userObj.Password,
		userObj.Email)

	if err != nil {
		panic(err)
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

func getReservationsForUser(userID int) []Reservation {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query(
		"SELECT r.id, r.start_time, r.end_time, r.car_id, r.garage_id "+
		"FROM reservations r "+
		"INNER JOIN cars on r.car_id = cars.id "+
		"INNER JOIN users on users.id = cars.user_id "+
		"WHERE users.id = $1",
		userID)
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
			Reservation{ID: id, StartTime: start, EndTime: end, CarID: carID, GarageID: garageID})
	}
	return reservations

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
			Reservation{ID: id, StartTime: start, EndTime: end, CarID: carID, GarageID: garageID})
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

		return Reservation{ID: id, StartTime: start, EndTime: end, CarID: carID, GarageID: garageID}
	}
	return Reservation{}
}

func getReservation(reservationID int) (Reservation) {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query(
		"SELECT * FROM reservations WHERE id = $1",
		reservationID)

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

		return Reservation{ID: id, StartTime: start, EndTime: end, CarID: carID, GarageID: garageID}
	}
	return Reservation{}
}

func getGarages() ([]Garage) {
	db := dbLogin()
	defer db.Close()

	rows, _ := db.Query("SELECT * FROM garages")
	var garages []Garage

	for rows.Next() {
		var name string
		var maxCars int
		var id int

		err := rows.Scan(&id, &name, &maxCars)
		if err != nil {
			panic(err)
		}
		garages = append(garages, Garage{ID: id, Name: name, MaxCars: maxCars})
	}
	return garages
}

func getGarage(garageID int) (Garage) {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query(
		"SELECT * FROM garages WHERE id = $1",
		garageID)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id int
		var name string
		var maxCars int

		scanErr := rows.Scan(&id, &name, &maxCars)
		if scanErr != nil {
			panic(err)
		}

		return Garage{ID: id, Name: name, MaxCars: maxCars}
	}
	return Garage{}
}

func getCars() ([]Car) {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query(
		"SELECT * FROM cars")

	var cars []Car
	for rows.Next() {
		var id int
		var user_id int
		var model string

		err = rows.Scan(&id, &user_id, &model)
		if err != nil {
			pnaic(err)
		}
		cars = append(cars, Car{ID: id, Model: model, UserID: user_id})
	}
	return cars
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

		userCars = append(userCars, Car{ID: id, Model: model})
	}
	return userCars
}

func getCar(carID int) (Car) {
	db := dbLogin()
	defer db.Close()

	rows, err := db.Query(
		"SELECT * FROM cars WHERE id = $1",
		carID)

	for rows.Next() {
		var id int
		var user_id int
		var model string

		err = rows.Scan(&id, &user_id, &model)
		if err != nil {
			panic(err)
		}
		return Car{ID: id, UserID: user_id, Model: model}
	}
	return Car{}
}

func getUserById(userID int) (User) {
	db := dbLogin()
	defer db.Close()

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

		return User{ID: id, Username: db_username, Email: email}
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

