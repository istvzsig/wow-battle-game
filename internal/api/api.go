package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/gorilla/mux"
	"github.com/istvzsig/wow-battle-game/internal/auth"
	"github.com/istvzsig/wow-battle-game/internal/entity"
	"github.com/istvzsig/wow-battle-game/internal/logger"
	"github.com/istvzsig/wow-battle-game/pkg/battle"
	"github.com/istvzsig/wow-battle-game/pkg/character"

	"google.golang.org/api/option"
)

type ApiServer struct {
	Address string
	Port    int
	Router  *mux.Router
	Db      *db.Client
	Logger  *log.Logger
}

func NewApiServer(address string, port int) *ApiServer {
	return &ApiServer{
		Address: address,
		Port:    port,
		Router:  mux.NewRouter(),
	}
}

func (s *ApiServer) Run() {
	startMsg := fmt.Sprintf("Running on http://%v:%v", s.Address, s.Port)
	s.Logger.Println(startMsg)
	log.Println(startMsg)
	err := http.ListenAndServe(":"+strconv.Itoa(s.Port), s.Router)
	if err != nil {
		s.Logger.Fatalf("Server failed to start: %v", err)
	}

}

func (s *ApiServer) InitRouter() {
	s.Router.HandleFunc("/account", HandleCreateAccount)

	s.Router.HandleFunc("/create", auth.AuthMiddleware(HandleCreateCharacter))
	s.Logger.Println("Registered protected endpoint: /create")

	s.Router.HandleFunc("/battle", auth.AuthMiddleware(HandleBattle))
	s.Logger.Println("Registered protected endpoint: /battle")

	log.Println("Router initialized.")
}

func (s *ApiServer) InitFirestore() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	googleKey := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	opt := option.WithCredentialsFile(googleKey)

	config := new(firebase.Config)
	config.DatabaseURL = os.Getenv("DATABASE_URL")

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	db, err := app.Database(ctx)
	if err != nil {
		log.Fatalf("Error getting Database client: %v", err)
		return
	}

	s.Db = db
	log.Println("Database initialized.")
}

func (s *ApiServer) InitLogger() {
	logger, err := logger.NewLogger("server.log")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
		return
	}

	s.Logger = logger

	log.Println("Logger initialized.")
}

func HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var acc entity.Entity = new(character.Character)
	acc.Create(w, r)
}

func HandleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var char entity.Entity = new(character.Character)
	char.Create(w, r)
}

func HandleBattle(w http.ResponseWriter, r *http.Request) {
	var battleResult = new(battle.BattleResult)
	battleResult.Create(w, r)
}
