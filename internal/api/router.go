package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/api/users/signup", SignUpHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/users/login", LoginHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/users/add_body_details", AddBodyDetailsHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/users/set_weight_goal", SetUserWeightGoalHandler).Methods(http.MethodPost)

	return router
}
