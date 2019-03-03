package main

import(
	"net/http"
	"encoding/json"
	"fmt"
)

func handleRaftJoinRequest(w http.ResponseWriter, req *http.Request) {
	m := map[string]string{}
	json.NewDecoder(req.Body).Decode(&m)

	remoteAddr,_ := m["addr"]
	nodeID,_ := m["id"]
	fmt.Println(remoteAddr, nodeID)

	s.Join(nodeID, remoteAddr)
}

func handleRaftDump(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(s.GetAll())
}
