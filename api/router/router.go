package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type cont interface {
	HelloWorld(w http.ResponseWriter, r *http.Request)
	GetContestans(w http.ResponseWriter, r *http.Request)
	GetSingleContestant(w http.ResponseWriter, r *http.Request)
}

var Resp = render.New()

func NewRouter(c cont) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router = router.PathPrefix("/api").Subrouter()

	router.HandleFunc("/", c.HelloWorld)
	router.HandleFunc("/contestants", c.GetContestans)
	router.HandleFunc("/contestants/{id}", c.GetSingleContestant)

	return router
}
