package api

import (
	"log"
	"net/http"

	"github.com/istvzsig/wow-battle-game/internal/account"
	"github.com/istvzsig/wow-battle-game/internal/entity"
	"github.com/istvzsig/wow-battle-game/pkg/battle"
)

func HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	var acc entity.Entity = new(account.Account)
	acc.Create(w, r)
}

func HandleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	// var char entity.Entity = new(character.Character)
	// char.Create(w, r, db)
}

func HandleBattle(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	var battleResult = new(battle.BattleResult)
	battleResult.Create(w, r)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	log.Println("Handling login...")
}
