package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Middlewares
	enableCORSMiddleware := alice.New(EnableCORS)
	authMiddleware := enableCORSMiddleware.Append(AuthMiddleWare)

	// Define routes
	router.Handle("/api/users/signup", enableCORSMiddleware.Then(http.HandlerFunc(SignUpHandler))).Methods(http.MethodPost)
	router.Handle("/api/users/login", enableCORSMiddleware.Then(http.HandlerFunc(LoginHandler))).Methods(http.MethodPost)

	router.Handle("/api/users/body_details/add", authMiddleware.Then(http.HandlerFunc(AddBodyDetailsHandler))).Methods(http.MethodPost)
	router.Handle("/api/users/body_details/exists", authMiddleware.Then(http.HandlerFunc(SetUserWeightGoalHandler))).Methods(http.MethodPost)
	router.Handle("/api/users/weight_goal/set", authMiddleware.Then(http.HandlerFunc(SetUserWeightGoalHandler))).Methods(http.MethodPost)

	router.Handle("/api/users/log/create", authMiddleware.Then(http.HandlerFunc(CreateCalorieLogHandler))).Methods(http.MethodPost)
	router.Handle("/api/users/log/get", authMiddleware.Then(http.HandlerFunc(GetCalorieLogsHandler))).Methods(http.MethodGet)
	router.Handle("/api/users/log/update", authMiddleware.Then(http.HandlerFunc(UpdateCalorieLogHandler))).Methods(http.MethodPut)
	router.Handle("/api/users/log/mark_status", authMiddleware.Then(http.HandlerFunc(MarkLoggingStatusHandler))).Methods(http.MethodPost)
	router.Handle("/api/users/log/delete", authMiddleware.Then(http.HandlerFunc(DeleteCalorieLogHandler))).Methods(http.MethodDelete)

	router.Handle("/api/users/net_caloric_balance/get", authMiddleware.Then(http.HandlerFunc(GetNetCaloricBalanceHandler))).Methods(http.MethodGet)

	return router
}
