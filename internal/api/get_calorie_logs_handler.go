package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type GetCaloricLogsHandlerResp struct {
	Logs *[]lib.UserCalorieLogs `json:"logs"`
}

func GetCalorieLogsHandler(w http.ResponseWriter, r *http.Request) {
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

	var userCalorieLogs *[]lib.UserCalorieLogs
	userCalorieLogs, err = lib.GetCalorieLogs(*userId)
	if err != nil {
		log.Info(
			"failed to get logs by user id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	data := &GetCaloricLogsHandlerResp{
		Logs: userCalorieLogs,
	}
	resp := Response{
		Code: http.StatusOK,
		Data: data,
	}
	json.NewEncoder(w).Encode(&resp)
}
