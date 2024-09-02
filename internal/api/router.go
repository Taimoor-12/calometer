package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/api/", GreetingHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/users/signup", SignUpHandler).Methods(http.MethodPost)

	return router
}
