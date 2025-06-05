package main

import (
	"encoding/json"
	"net/http"

	"github.com/harshadmanglani/poseidon/clients"
	"github.com/harshadmanglani/poseidon/config"
	"github.com/harshadmanglani/poseidon/db"
	"github.com/harshadmanglani/poseidon/workflows"
)

type Context struct {
	Service   string `json:"service"`
	Type      string `json:"type"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

type InvokeRequest struct {
	Callback string     `json:"callback"`
	ID       string     `json:"id"`
	Context  []*Context `json:"context"`
}

func handleInvoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InvokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func StartServer() {
	http.HandleFunc("/invoke", handleInvoke)
	http.ListenAndServe(":8080", nil)
}

func main() {

	config.Init()

	clients.Init()

	db.Init()

	workflows.Init()

	StartServer()
}
