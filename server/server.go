package server

import (
	"log"

	"github.com/istvzsig/wow-battle-game/api"
)

var API *api.ApiServer

func InitAPI(host string, port int) {
	API = api.NewApiServer(host, port)
	API.InitLogger()
	API.InitRouter()
	API.InitFirestore()
}

func RunAPI() {
	if API != nil {
		API.Run()
	} else {
		log.Fatal("API server is not initialized")
	}
}
