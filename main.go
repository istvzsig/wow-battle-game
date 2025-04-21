package main

import (
	"github.com/istvzsig/wow-battle-game/internal/api"

	"github.com/istvzsig/wow-battle-game/internal/server"
)

func main() {
	var api = api.NewApiServer("localhost", 8080)

	server.Init(api)
	server.Run(api)
}
