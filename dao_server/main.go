package main

import (
	"log"
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"os/exec"
	"strings"
	"time"
	"github.com/gorilla/mux"
	"github.com/otoolep/hraftd/store"
)

const (
	retainSnapshotCount = 2
	raftTimeout			= 10 * time.Second
)

// Store is the interface Raft-backed key-value stores must implement.
type Store interface {
	// Get returns the value for the given key.
	Get(key string) (string, error)

	// Set sets the value for the given key, via distributed consensus.
	Set(key, value string) error

	// Delete removes the given key, via distributed consensus.
	Delete(key string) error

	// Join joins the node, identitifed by nodeID and reachable at addr, to the cluster.
	Join(nodeID string, addr string) error
}

//Use "go run *.go" to run the program
// Main function
func main() {
	// Init router
	router := mux.NewRouter()

	// Reservation route handles & endpoints
	router.HandleFunc("/reservations", GetReservations).Methods("GET")
	router.HandleFunc("/reservations/", GetReservations).Methods("GET")
	router.HandleFunc("/reservations/{id}", GetReservation).Methods("GET")
	router.HandleFunc("/reservations/", CreateReservation).Methods("POST")
	router.HandleFunc("/reservations", CreateReservation).Methods("POST")
	//router.HandleFunc("/reservations/{id}", UpdateReservation).Methods("PUT")
	router.HandleFunc("/reservations/{id}", DeleteReservation).Methods("DELETE")
	router.HandleFunc("/reservations-by-user/{id}", GetReservationsByUser).Methods("GET")

	// Users route handles & endpoints
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users/", GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/", CreateUser).Methods("POST")
	//router.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	// Garages route handles & endpoints
	router.HandleFunc("/garages", GetGarages).Methods("GET")
	router.HandleFunc("/garages/", GetGarages).Methods("GET")
	router.HandleFunc("/garages/{id}", GetGarage).Methods("GET")
	router.HandleFunc("/garages/", CreateGarage).Methods("POST")
	router.HandleFunc("/garages", CreateGarage).Methods("POST")
	//router.HandleFunc("/garages/{id}", UpdateGarage).Methods("PUT")
	router.HandleFunc("/garages/{id}", DeleteGarage).Methods("DELETE")

	// Cars route handles & endpoints
	router.HandleFunc("/cars", GetCars).Methods("GET")
	router.HandleFunc("/cars/", GetCars).Methods("GET")
	router.HandleFunc("/cars/{id}", GetCar).Methods("GET")
	router.HandleFunc("/cars", CreateCar).Methods("POST")
	router.HandleFunc("/cars", CreateCar).Methods("POST")
	router.HandleFunc("/get-cars-by-user/{id}", GetCarsByUser).Methods("GET")
	//router.HandleFunc("/cars/{id}", UpdateCar).Methods("PUT")
	router.HandleFunc("/cars/{id}", DeleteCar).Methods("DELETE")

	// History route handler & endpoint
	router.HandleFunc("/history/", GetHistory).Methods("GET")

	//router.HandleFunc("/join/", handleRaftJoinRequest).Methods("POST")

	// Start server
	go func() {
		log.Fatal(http.ListenAndServe(":8888", router))
	}()

	// Get LAN IP (private IP in GCP console)
	cmd := "ifconfig | grep 'inet 10' | awk '{print $2}'"
	out, _ := exec.Command("bash", "-c", cmd).Output()
	localIP := string(out)
	raftIP := strings.Replace(localIP+":12000", "\n", "", -1)

	// Directory where Raft logs and ledgers will be stored
	raftNodeStoreDir := "node.secret"
	os.MkdirAll(raftNodeStoreDir, 0700)

	s := store.New(false)
	s.RaftDir = raftNodeStoreDir
	s.RaftBind = raftIP

	// Hardcoded for now, will probably configure via env for fresh Raft startup
	isFirstNode := true

	if err := s.Open(isFirstNode, "node0"); err != nil {
		log.Fatalf("Failed to opeen store: %s", err.Error())
	}

	// Block on a channel that waits for a SIGKILL, can kill via Ctrl+C
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	fmt.Println("Shutting down node")
}
