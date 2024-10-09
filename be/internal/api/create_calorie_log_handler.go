package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type CreateCalorieLogReq struct {
	LogDate time.Time `json:"log_date,omitempty"`
}

func CreateCalorieLogHandler(w http.ResponseWriter, r *http.Request) {
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

	var req CreateCalorieLogReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var logDate string
	if req.LogDate.IsZero() {
		logDate = time.Now().Format("2006-01-02")
	} else {
		// Check if the logDate is in the future
		if req.LogDate.After(time.Now()) {
			http.Error(w, "Log date cannot be a future date", http.StatusBadRequest)
			return
		}

		logDate = req.LogDate.Format("2006-01-02")
	}

	exists, err := lib.DoesLogExistForTheDay(*userId, logDate)
	if err != nil {
		log.Info(
			"failed to determine user log's existence",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if *exists {
		http.Error(w, "Log already exists for this day", http.StatusConflict)
		return
	}

	bmr, err := lib.GetUserBmr(*userId)
	if err != nil {
		log.Info(
			"failed to get user bmr by id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := lib.CreateUserLog(*userId, *bmr, logDate); err != nil {
		log.Info(
			"failed to create user log by id",
			zap.String("userId", userId.String()),
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
