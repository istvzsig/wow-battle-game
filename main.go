package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Player represents the main character controlled by the user.
type Player struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Level int    `json:"level"`
	HP    int    `json:"hp"` // Hit Points
	AP    int    `json:"ap"` // Attack Power
}

// NewPlayer returns a pointer to a new empty Player instance.
func NewPlayer() *Player {
	return &Player{}
}

// Monster represents a basic enemy NPC for battles.
type Monster struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Level int    `json:"level"`
	HP    int    `json:"hp"` // Hit Points
	AP    int    `json:"ap"` // Attack Power
}

func NewMonster(name, cls string, level, hp, ap int) *Monster {
	return &Monster{
		Name:  name,
		Class: cls,
		Level: level,
		HP:    hp,
		AP:    ap,
	}
}

// BattleResult stores the outcome of a single battle.
type BattleResult struct {
	PlayerHP  int    `json:"player_hp"`  // Remaining HP of the player
	MonsterHP int    `json:"monster_hp"` // Remaining HP of the monster
	Winner    string `json:"winner"`     // "Player" or "Monster"
}

// CreateCharacterResponse defines the API response when a new character is created.
type CreateCharacterResponse struct {
	Status    int     `json:"status"`    // HTTP status code
	Character *Player `json:"character"` // Newly created player data
}

type BattleResultResponse struct {
	Status int           `json:"status"` // HTTP status code
	Result *BattleResult `json:"battleResult"`
}

// main sets up HTTP routes and starts the server.
func main() {
	mux := http.NewServeMux()

	// Register endpoints
	mux.HandleFunc("/create", handleCreate)
	mux.HandleFunc("/battle", handleBattle)

	log.Println("Running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// handleCreate handles the character creation request from the frontend.
func handleCreate(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	// Handle OPTIONS (CORS preflight)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body into Player struct
	var p = NewPlayer()
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Initialize base stats for the new character
	p.Level = 1
	p.HP = 100
	p.AP = 15

	// Send response back with created character
	json.NewEncoder(w).Encode(CreateCharacterResponse{
		Status:    http.StatusCreated,
		Character: p,
	})
}

// handleBattle processes a fight between a Player and a Monster.
func handleBattle(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	// Handle OPTIONS (CORS preflight)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode player from request
	var p Player
	_ = json.NewDecoder(r.Body).Decode(&p)

	// Create a default monster
	m := NewMonster("Board", "Creature", p.Level, 100, 20)
	rand.Seed(time.Now().UnixNano())

	// Simulate turn-based battle
	for p.HP > 0 && m.HP > 0 {
		m.HP -= rand.Intn(p.AP)
		if m.HP <= 0 {
			break
		}
		p.HP -= rand.Intn(m.AP)
	}

	// Determine winner
	result := BattleResult{
		PlayerHP:  p.HP,
		MonsterHP: m.HP,
		Winner:    "Monster",
	}
	if p.HP > 0 {
		result.Winner = "Player"
	}

	// Send result as JSON response
	json.NewEncoder(w).Encode(BattleResultResponse{
		Status: http.StatusOK,
		Result: &result,
	})
}

// enableCORS sets the necessary headers to allow cross-origin requests from the frontend.
func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // Adjust to specific origin for more security
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
