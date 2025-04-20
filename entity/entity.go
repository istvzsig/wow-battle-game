package entity

import "net/http"

type Entity interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, id any)
	Update(w http.ResponseWriter, r *http.Request, id any)
	Delete(w http.ResponseWriter, r *http.Request, id any)
}
