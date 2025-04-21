package entity

import (
	"net/http"

	"firebase.google.com/go/db"
)

type Entity interface {
	Create(w http.ResponseWriter, r *http.Request, db *db.Client)
	Get(w http.ResponseWriter, r *http.Request, id any)
	Update(w http.ResponseWriter, r *http.Request, id any)
	Delete(w http.ResponseWriter, r *http.Request, id any)
}
