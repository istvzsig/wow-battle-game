package server

import (
	"log"
	"os"

	"github.com/istvzsig/wow-battle-game/internal/api"
	"github.com/joho/godotenv"
)

func Run(api *api.ApiServer) {
	if api != nil {
		api.Run()
	} else {
		log.Fatal("API server is not initialized")
	}
}

func Init(api *api.ApiServer) {
	LoadEnv(api)
	api.InitLogger("config.log")
	api.InitRouter()
	api.InitFirestore(
		os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		os.Getenv("DATABASE_URL"),
	)
}

func LoadEnv(api *api.ApiServer) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		api.Logger.Fatal("Error loading .env file")
	}
}
