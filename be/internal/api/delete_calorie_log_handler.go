package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type DeleteCalorieLogReq struct {
	LogDate time.Time `json:"log_date"`
}

func DeleteCalorieLogHandler(w http.ResponseWriter, r *http.Request) {
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

	var req DeleteCalorieLogReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var logDate string
	// Check if the logDate is in the future
	if req.LogDate.After(time.Now()) {
		http.Error(w, "Log date cannot be a future date", http.StatusBadRequest)
		return
	}

	logDate = req.LogDate.Format("2006-01-02")

	exists, err := lib.DoesLogExistForTheDay(*userId, logDate)
	if err != nil {
		log.Info(
			"failed to determine user log's existence by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !*exists {
		http.Error(w, "No log exists for this day", http.StatusConflict)
		return
	}

	logId, err := lib.GetCalorieLogId(*userId, logDate)
	if err != nil {
		log.Info(
			"failed to get logId by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := lib.DeleteCaloricBalanceByLogId(*logId); err != nil {
		log.Info(
			"failed to delete caloric balance by log id",
			zap.String("logId", logId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := lib.DeleteCalorieLog(*userId, logDate); err != nil {
		log.Info(
			"failed to delete log by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{
		Code: http.StatusOK,
	}
	json.NewEncoder(w).Encode(&resp)
}
