package main

import (
	"log"
	"net/http"
)

type Entity interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, id any)
	Update(w http.ResponseWriter, r *http.Request, id any)
	Delete(w http.ResponseWriter, r *http.Request, id any)
}

// main sets up HTTP routes and starts the server.
func main() {
	mux := http.NewServeMux()

	// Register endpoints
	mux.HandleFunc("/account", handleCreateAccount)
	mux.HandleFunc("/create", handleCreateCharacter)
	mux.HandleFunc("/battle", handleBattle)

	log.Println("Running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// handleCreateCharacterAccount handles the character creation request from the frontend.
func handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc Entity = NewAccount()
	acc.Create(w, r)
}

// handleCreateCharacter handles the character creation request from the frontend.
func handleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var char Entity = NewCharacter()
	char.Create(w, r)

}

// handleBattle processes a fight between a Player and a Monster.
func handleBattle(w http.ResponseWriter, r *http.Request) {
	battle := new(BattleResult)
	battle.Create(w, r)
}

// enableCORS sets the necessary headers to allow cross-origin requests from the frontend.
func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
