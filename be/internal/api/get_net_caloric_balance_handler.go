package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"math"
	"net/http"

	"go.uber.org/zap"
)

type GetNetCaloricBalanceHandlerResp struct {
	NetCaloricBalance float64 `json:"net_caloric_balance"`
}

func GetNetCaloricBalanceHandler(w http.ResponseWriter, r *http.Request) {
	resp := Response{}
	resp.Code = make(map[int]string)

	// Retrieve the token from the context
	tokenStr, ok := r.Context().Value(TokenContextKey).(string)
	if !ok {
		log.Info(
			"token not found in context",
		)

		// Token is not present in context
		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	userId, err := lib.ExtractUserIdFromToken(tokenStr)
	if err != nil {
		log.Info(
			"failed to get user id by token",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	netCaloricBalance, err := lib.GetNetCaloricBalance(*userId)
	if err != nil {
		log.Info(
			"failed to get net caloric balance by id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	weightGoal, err := lib.GetUserWeightGoalById(*userId)
	if err != nil {
		log.Info(
			"failed to get user weight goal by id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if *netCaloricBalance < 0 && *weightGoal == "G" {
		*netCaloricBalance = math.Abs(*netCaloricBalance)
	} else if *netCaloricBalance > 0 && *weightGoal == "G" {
		*netCaloricBalance = -*netCaloricBalance
	}

	data := &GetNetCaloricBalanceHandlerResp{
		NetCaloricBalance: *netCaloricBalance,
	}

	resp.Code[http.StatusOK] = "OK"
	resp.Data = data

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}
