package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Host struct {
	Name string
}

func (h *Host) Healtz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"host": h.Name, "status": "ready"})
}

func main() {
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal("unable to get hostname")
	}
	host := Host{hostName}

	http.HandleFunc("/healtz", host.Healtz)
	if http.ListenAndServe(":8080", nil) != nil {
		log.Fatal("unable to start http server")
	}
}
