package main

import (
	"encoding/json"
	"net/http"
)

// Account represents the auth for user.
type Account struct {
	Name     string
	Email    string
	Password string
}

// CreateAccountResponse defines the API response when a new account is created.
type CreateAccountResponse struct {
	Status  int `json:"status"`  // HTTP status code
	Account any `json:"account"` // Newly created player data
}

func NewAccount() *Account {
	return &Account{}
}

func (acc *Account) Create(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	// Handle OPTIONS (CORS preflight)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body into Player struct
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Send response back with created character
	json.NewEncoder(w).Encode(CreateAccountResponse{
		Status:  http.StatusCreated,
		Account: acc,
	})
}

func (acc *Account) Get(w http.ResponseWriter, r *http.Request, id any) {

}

func (acc *Account) Update(w http.ResponseWriter, r *http.Request, id any) {

}

func (acc *Account) Delete(w http.ResponseWriter, r *http.Request, id any) {

}
