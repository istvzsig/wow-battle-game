package main

import (
	"log"

	"github.com/istvzsig/wow-battle-game/server"
	"github.com/joho/godotenv"
)

func main() {

	loadEnv()
	server.InitAPI("localhost", 8080)
	server.RunAPI()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		server.API.Logger.Fatal("Error loading .env file")
	}
}
