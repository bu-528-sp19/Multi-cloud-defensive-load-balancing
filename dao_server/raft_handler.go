package main

import(
	"net/http"
	"encoding/json"
)

func handleRaftJoinRequest(w http.ResponseWriter, req *http.Request) {
	m := map[string]string{}
	json.NewDecoder(req.Body).Decode(&m)

	remoteAddr,_ := m["addr"]
	nodeID,_ := m["id"]

	s.Join(nodeID, remoteAddr)
}
