package main

import (
	"encoding/json"
	"net/http"

	"github.com/harshadmanglani/poseidon/clients"
	"github.com/harshadmanglani/poseidon/config"
	"github.com/harshadmanglani/poseidon/db"
	"github.com/harshadmanglani/poseidon/utils"
	"github.com/harshadmanglani/poseidon/workflows"
)

type InvokeRequest struct {
	ID      string                `json:"id"`
	Context workflows.ContextData `json:"context"`
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

	analysis, err := workflows.Invoke(req.ID, req.Context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}

func StartServer() {
	utils.Sugar.Infof("Starting Poseidon server on port 8080...")
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
