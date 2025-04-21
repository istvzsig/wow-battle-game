package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"

	"github.com/istvzsig/wow-battle-game/internal/account"
	"github.com/istvzsig/wow-battle-game/internal/auth"
	"github.com/istvzsig/wow-battle-game/internal/entity"
	"github.com/istvzsig/wow-battle-game/internal/logger"
	"github.com/istvzsig/wow-battle-game/pkg/battle"
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
	s.Router.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		HandleCreateAccount(w, r, s.Db)
	})

	s.Router.HandleFunc("/login", HandleLogin)

	s.Router.HandleFunc("/create", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		HandleCreateCharacter(w, r)
	}))
	s.Logger.Println("Registered protected endpoint: /create")

	s.Router.HandleFunc("/battle", auth.AuthMiddleware(HandleBattle))
	s.Logger.Println("Registered protected endpoint: /battle")

	log.Println("Router initialized.")
}

func (s *ApiServer) InitFirestore(key, dbUrl string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opt := option.WithCredentialsFile(key)

	config := new(firebase.Config)
	config.DatabaseURL = dbUrl

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

func HandleCreateAccount(w http.ResponseWriter, r *http.Request, db *db.Client) {
	enableCORS(&w)
	var acc entity.Entity = new(account.Account)
	acc.Create(w, r, db)
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
