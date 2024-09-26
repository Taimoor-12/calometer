package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"math"
	"net/http"

	"go.uber.org/zap"
)

func GetNetCaloricBalanceHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the token from the context
	tokenStr, ok := r.Context().Value(TokenContextKey).(string)
	if !ok {
		log.Info(
			"token not found in context",
		)

		// Token is not present in context
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userId, err := lib.ExtractUserIdFromToken(tokenStr)
	if err != nil {
		log.Info(
			"failed to get user id by token",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	netCaloricBalance, err := lib.GetNetCaloricBalance(*userId)
	if err != nil {
		log.Info(
			"failed to get net caloric balance by id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	weightGoal, err := lib.GetUserWeightGoalById(*userId)
	if err != nil {
		log.Info(
			"failed to get user weight goal by id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if *netCaloricBalance < 0 && *weightGoal == "G" {
		*netCaloricBalance = math.Abs(*netCaloricBalance)
	} else if *netCaloricBalance > 0 && *weightGoal == "G" {
		*netCaloricBalance = -*netCaloricBalance
	}

	resp := Response{
		Code: http.StatusOK,
		Data: *netCaloricBalance,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
