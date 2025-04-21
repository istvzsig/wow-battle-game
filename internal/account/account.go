package account

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	firebase "firebase.google.com/go/db"
	"github.com/google/uuid"

	"github.com/istvzsig/wow-battle-game/internal/db"
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
	Status   int    `json:"status"`
	Msg      string `json:"msg"`
	*Account `json:"account,omitempty"`
}

// var mu *sync.Mutex

// Create Account
func (acc *Account) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger, err := logger.NewLogger("account.log")
	if err != nil {
		log.Fatal("Failed to create account logger:", err)
		return
	}

	if r.Method == http.MethodOptions {
		WithJSONResponse(w, http.StatusAccepted)
		logger.Printf("CORS preflight request from %s\n", r.RemoteAddr)
		return
	}

	if r.Method != http.MethodPost {
		WithJSONResponse(w, http.StatusMethodNotAllowed)
		logger.Fatalf("Rejected non-POST request from %s: %s\n", r.RemoteAddr, r.Method)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		WithJSONResponse(w, http.StatusInternalServerError)
		logger.Fatalf("JSON decode error from %s: %v\n", r.RemoteAddr, err)
		return
	}

	acc.ID = uuid.New().String()
	acc.CreatedAt = time.Now().Unix()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
	if err != nil {
		WithJSONResponse(w, http.StatusInternalServerError, "Failed to hash password.")
		logger.Fatalf("Password hashing error from %s: %v\n", r.RemoteAddr, err)
		return
	}
	acc.Password = string(hashedPassword)

	ref := db.FirestoreClient.NewRef("users").Child(acc.ID)

	// Check if the account already exists
	accList, err := GetUsers(ctx, db.FirestoreClient.NewRef("users"))
	if err != nil {
		logger.Fatalf("Cannot find accounts list %s: %v\n", r.RemoteAddr, err)
		return
	}

	// Check for duplicate accounts based on a unique field (e.g., email)
	for _, existingAcc := range accList.Items {
		if existingAcc.Email == acc.Email || existingAcc.Name == acc.Name || existingAcc.Password == acc.Password { // Assuming Email is a field in Account
			WithJSONResponse(w, http.StatusConflict, "Account already exists.")
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

	WithJSONResponse(w, http.StatusCreated, "Account created")
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
func GetUsers(ctx context.Context, ref *firebase.Ref) (*AccountList, error) {
	accountList := new(AccountList)

	if err := ref.Get(ctx, &accountList.Items); err != nil {
		return nil, err
	}

	return accountList, nil
}

func WithJSONResponse(w http.ResponseWriter, status int, args ...any) {
	w.WriteHeader(status)

	response := AccountCreateResponse{
		Status: status,
	}

	// Use a switch statement to handle different types of arguments
	switch len(args) {
	case 1:
		// If there's one argument, check its type
		switch v := args[0].(type) {
		case string:
			response.Msg = v // If it's a string, set it as the message
		case Account: // Assuming Account is a defined struct
			response.Account = &v // If it's an Account, set it
		default:
			response.Msg = "Unknown error" // Fallback message
		}
	case 2:
		// If there are two arguments, handle them accordingly
		if msg, ok := args[0].(string); ok {
			response.Msg = msg
		}
		response.Account = args[1].(*Account)
	}

	// Encode the response as JSON
	json.NewEncoder(w).Encode(response)
}
