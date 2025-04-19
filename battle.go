package main

import "net/http"

// BattleResult stores the outcome of a single battle.
type BattleResult struct {
	PlayerHP  int    `json:"player_hp"`  // Remaining HP of the player
	MonsterHP int    `json:"monster_hp"` // Remaining HP of the monster
	Winner    string `json:"winner"`     // "Player" or "Monster"
}

type BattleResultResponse struct {
	Status int          `json:"status"` // HTTP status code
	Result BattleResult `json:"battleResult"`
}

func (b *BattleResult) Create(w http.ResponseWriter, r *http.Request) {

}
