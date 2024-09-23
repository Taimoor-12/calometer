package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Middlewares
	authMiddleware := alice.New(AuthMiddleWare)

	// Define routes
	router.HandleFunc("/api/users/signup", SignUpHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/users/login", LoginHandler).Methods(http.MethodPost)

	router.Handle("/api/users/add_body_details", authMiddleware.Then(http.HandlerFunc(AddBodyDetailsHandler))).Methods(http.MethodPost)
	router.Handle("/api/users/set_weight_goal", authMiddleware.Then(http.HandlerFunc(SetUserWeightGoalHandler))).Methods(http.MethodPost)

	router.Handle("/api/users/log/create", authMiddleware.Then(http.HandlerFunc(CreateCalorieLogHandler))).Methods(http.MethodPost)
	router.Handle("/api/users/log/get", authMiddleware.Then(http.HandlerFunc(GetCalorieLogs))).Methods(http.MethodGet)

	router.Handle("/api/users/log_calories/consumed", authMiddleware.Then(http.HandlerFunc(LogCaloriesConsumedHandler))).Methods(http.MethodPost)
	router.Handle("/api/users/log_calories/burnt", authMiddleware.Then(http.HandlerFunc(LogCaloriesBurntHandler))).Methods(http.MethodPost)
	router.Handle("/api/users/log_calories/mark_status", authMiddleware.Then(http.HandlerFunc(MarkLoggingStatusHandler))).Methods(http.MethodPost)

	return router
}
