package auth

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/istvzsig/wow-battle-game/internal/logger"
)

var jwtKey = []byte(os.Getenv("super_secret_key"))

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger, err := logger.NewLogger("auth.log")
		if err != nil {
			log.Println("Could not create logger:", err)
			// TODO: Add fallback logging to default
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 8 {
			if logger != nil {
				logger.Println("Missing or malformed token from", r.RemoteAddr)
			}
			http.Error(w, "Missing or malformed token", http.StatusUnauthorized)
			return
		}

		tokenStr := authHeader[len("Bearer "):]
		_, err = ValidateJWT(tokenStr)
		if err != nil {
			if logger != nil {
				logger.Printf("Invalid token from %s: %v\n", r.RemoteAddr, err)
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if logger != nil {
			logger.Printf("Authorized request from %s\n", r.RemoteAddr)
		}

		next(w, r)
	}
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
