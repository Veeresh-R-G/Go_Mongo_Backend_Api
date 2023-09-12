package router

import (
	"github.com/Veeresh-R-G/mongoapi/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/api/movies", controller.GetAllMovies).Methods("GET")
	router.HandleFunc("/api/movies", controller.CreateMovie).Methods("POST")
	router.HandleFunc("/api/movies/{id}", controller.MarkedAsWatched).Methods("PUT")

	// Can add more routes based on controllers
	return router
}
