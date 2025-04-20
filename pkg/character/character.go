package character

import (
	"encoding/json"
	"net/http"
)

// Character represents the character controlled by the user.
type Character struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Level int    `json:"level"`
	HP    int    `json:"hp"` // Hit Points
	AP    int    `json:"ap"` // Attack Power
}

// CreateCharacterResponse defines the API response when a new character is created.
type CreateCharacterResponse struct {
	Status    int        `json:"status"`    // HTTP status code
	Character *Character `json:"character"` // Newly created player data
}

// NewPlayer returns a pointer to a new empty Player instance.
func NewCharacter() *Character {
	return &Character{}
}

func (char *Character) Create(w http.ResponseWriter, r *http.Request) {
	// enableCORS(&w)

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
	if err := json.NewDecoder(r.Body).Decode(&char); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Send response back with created character
	json.NewEncoder(w).Encode(CreateCharacterResponse{
		Status:    http.StatusCreated,
		Character: char,
	})
}

func (char *Character) Get(w http.ResponseWriter, r *http.Request, id any) {

}

func (char *Character) Update(w http.ResponseWriter, r *http.Request, id any) {

}

func (char *Character) Delete(w http.ResponseWriter, r *http.Request, id any) {

}
