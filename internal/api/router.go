package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	logCaloriesConsumedMiddleware := alice.New(AuthMiddleWare, SetInitialTDEEMiddleware)

	// Define routes
	router.HandleFunc("/api/users/signup", SignUpHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/users/login", LoginHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/users/add_body_details", AddBodyDetailsHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/users/set_weight_goal", SetUserWeightGoalHandler).Methods(http.MethodPost)

	router.Handle("/api/users/log_calories/consumed", logCaloriesConsumedMiddleware.Then(http.HandlerFunc(LogCaloriesConsumedHandler))).Methods(http.MethodPost)
	return router
}
