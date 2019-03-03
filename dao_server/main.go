package main

import (
	"log"
	"bytes"
	"encoding/json"
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

var s *store.Store

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
	router.HandleFunc("/cars/", CreateCar).Methods("POST")
	router.HandleFunc("/get-cars-by-user/{id}", GetCarsByUser).Methods("GET")
	//router.HandleFunc("/cars/{id}", UpdateCar).Methods("PUT")
	router.HandleFunc("/cars/{id}", DeleteCar).Methods("DELETE")

	// History route handler & endpoint
	router.HandleFunc("/history/", GetHistory).Methods("GET")

	router.HandleFunc("/join", handleRaftJoinRequest).Methods("POST")
	router.HandleFunc("/join/", handleRaftJoinRequest).Methods("POST")
	router.HandleFunc("/raft-dump", handleRaftDump).Methods("GET")
	router.HandleFunc("/raft-dump/", handleRaftDump).Methods("GET")

	// Get LAN IP (private IP in GCP console)
	cmd := "ifconfig | grep 'inet 10' | awk '{print $2}'"
	out, _ := exec.Command("bash", "-c", cmd).Output()
	localIP := string(out)
	raftIP := strings.Replace(localIP+":12000", "\n", "", -1)

	// Directory where Raft logs and ledgers will be stored
	raftNodeStoreDir := "node.secret"
	os.MkdirAll(raftNodeStoreDir, 0700)

	s = store.New(false)
	s.RaftDir = raftNodeStoreDir
	s.RaftBind = raftIP

	// Hardcoded for now, will probably configure via env for fresh Raft startup
	isLeader := os.Getenv("LEADER")
	var isFirstNode bool
	if (isLeader == "true") {
		isFirstNode = true
	} else {
		isFirstNode = false
	}

	nodeID := os.Getenv("RAFT_NODE_ID")

	if err := s.Open(isFirstNode, nodeID); err != nil {
		log.Fatalf("Failed to opeen store: %s", err.Error())
	}

	if (!isFirstNode) {
		ipCmd := os.Getenv("EXTERNAL_IP_QUERY")
		ip, _ := exec.Command("bash", "-c", ipCmd).Output()
		externalIP := string(ip)
		externalRaftIP := strings.Replace(externalIP+":12000", "\n", "", -1)
		leaderIP := os.Getenv("LEADER_IP")
		b, err := json.Marshal(map[string]string{"addr": externalRaftIP, "id":nodeID})

		if err != nil {
			panic(err)
		}

		_,_ = http.Post("http://"+leaderIP+"/join", "application-type/json", bytes.NewReader(b))
	}

	go func() {
		log.Fatal(http.ListenAndServe(":8888", router))
	}()

	// Block on a channel that waits for a SIGKILL, can kill via Ctrl+C
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	fmt.Println("Shutting down node")
}
