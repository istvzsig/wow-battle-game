package account

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"firebase.google.com/go/db"
	"github.com/google/uuid"

	"github.com/istvzsig/wow-battle-game/internal/logger"
)

type Account struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"createdAt"`
}

type AccountList struct {
	Items map[string]Account `json:"items"`
}

type AccountCreateResponse struct {
	Status   int `json:"status"`
	*Account `json:"account"`
}

// var mu *sync.Mutex

// Create Account
func (acc *Account) Create(w http.ResponseWriter, r *http.Request, db *db.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger, err := logger.NewLogger("account.log")
	if err != nil {
		log.Fatal("Failed to create account logger:", err)
		return
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusAccepted)
		logger.Printf("CORS preflight request from %s\n", r.RemoteAddr)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		logger.Fatalf("Rejected non-POST request from %s: %s\n", r.RemoteAddr, r.Method)
		log.Fatalf("Rejected non-POST request from %s: %s\n", r.RemoteAddr, r.Method)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		logger.Fatalf("JSON decode error from %s: %v\n", r.RemoteAddr, err)
		return
	}

	acc.ID = uuid.New().String()
	acc.CreatedAt = time.Now().Unix()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		logger.Fatalf("Password hashing error from %s: %v\n", r.RemoteAddr, err)
		return
	}
	acc.Password = string(hashedPassword)

	ref := db.NewRef("users").Child(acc.ID)

	// Check if the user already exists
	accList, err := GetUsers(ctx, db.NewRef("users"))
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		logger.Fatalf("Cannot find accounts list %s: %v\n", r.RemoteAddr, err)
		return
	}

	// Check for duplicate accounts based on a unique field (e.g., email)
	for _, existingAcc := range accList.Items {
		if existingAcc.Email == acc.Email { // Assuming Email is a field in Account
			http.Error(w, "User already exists", http.StatusConflict)
			logger.Printf("Duplicate account attempt from %s: %s\n", r.RemoteAddr, acc.Email)
			return
		}
	}

	// Persist the new account at the reference created using the ID
	if err := ref.Set(ctx, acc); err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		logger.Fatalf("Failed to create account in DB %s: %v\n", r.RemoteAddr, err)
		return
	}

	if err := json.NewEncoder(w).Encode(AccountCreateResponse{
		Status:  http.StatusCreated,
		Account: acc,
	}); err != nil {
		log.Printf("Failed to encode account response: %v", err)
		if logger != nil {
			logger.Printf("JSON encode error for %s: %v\n", r.RemoteAddr, err)
		}
	}
	logger.Printf("Account created with id %v", acc.ID)
}

// Get Account by ID
func (acc *Account) Get(w http.ResponseWriter, r *http.Request, id any) {

}

// Update Account by ID
func (acc *Account) Update(w http.ResponseWriter, r *http.Request, id any) {

}

// Flag Deleted the Account by ID
func (acc *Account) Delete(w http.ResponseWriter, r *http.Request, id any) {

}

// Get users list from db
func GetUsers(ctx context.Context, ref *db.Ref) (*AccountList, error) {
	accountList := new(AccountList)

	if err := ref.Get(ctx, &accountList.Items); err != nil {
		return nil, err
	}

	return accountList, nil
}
