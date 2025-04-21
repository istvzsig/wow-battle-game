package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"firebase.google.com/go/db"
	"github.com/gorilla/mux"

	"github.com/istvzsig/wow-battle-game/internal/auth"
	"github.com/istvzsig/wow-battle-game/internal/logger"
	"github.com/istvzsig/wow-battle-game/pkg/utils"
)

type ApiServer struct {
	Address string
	Port    string
	Router  *mux.Router
	Store   *db.Client
	Logger  *log.Logger
}

func NewApiServer() *ApiServer {
	utils.LoadEnv()
	return &ApiServer{
		Address: os.Getenv("BACKEND_URL"),
		Port:    os.Getenv("BACKEND_PORT"),
		Router:  mux.NewRouter(),
	}
}

func (s *ApiServer) Run() {
	startMsg := fmt.Sprintf("Running on http://%v:%v", s.Address, s.Port)
	s.Logger.Println(startMsg)
	log.Println(startMsg)
	err := http.ListenAndServe(":"+s.Port, s.Router)
	if err != nil {
		s.Logger.Fatalf("Server failed to start: %v", err)
	}

}

func (s *ApiServer) InitRouter() {
	s.Router.HandleFunc("/account", HandleCreateAccount)

	s.Router.HandleFunc("/login", HandleLogin)

	s.Router.HandleFunc("/create", auth.AuthMiddleware(HandleCreateCharacter))
	s.Logger.Println("Registered protected endpoint: /create")

	s.Router.HandleFunc("/battle", auth.AuthMiddleware(HandleBattle))
	s.Logger.Println("Registered protected endpoint: /battle")

	log.Println("Router initialized.")
}

func (s *ApiServer) InitLogger(fileName string) {
	logger, err := logger.NewLogger(fileName)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
		return
	}

	s.Logger = logger
	log.Println("Logger initialized.")
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Content-Type", "application/json")
}
