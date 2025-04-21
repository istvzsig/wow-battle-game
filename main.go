package main

import (
	"github.com/istvzsig/wow-battle-game/internal/api"
	"github.com/istvzsig/wow-battle-game/internal/db"
)

func main() {
	server := api.NewApiServer()

	db.InitFirestore()

	server.InitLogger("server.log")
	server.InitRouter()
	server.Run()
}
