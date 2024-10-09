package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type UpdateCalorieLogReq struct {
	CaloriesConsumed float64   `json:"calories_consumed"`
	CaloriesBurnt    float64   `json:"calories_burnt"`
	LogDate          time.Time `json:"log_date"`
}

func UpdateCalorieLogHandler(w http.ResponseWriter, r *http.Request) {
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

	var req UpdateCalorieLogReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	logDate := req.LogDate.Format("2006-01-02")

	logStatus, err := lib.CheckLogStatusByIdAndDate(*userId, logDate)
	if err != nil {
		log.Info(
			"failed to check log status by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if *logStatus == "D" {
		http.Error(w, "Log is already completed", http.StatusConflict)
		return
	}

	if req.CaloriesBurnt != 0.00 {
		currValue, err := lib.FetchCaloriesBurntForTheDay(*userId, logDate)
		if err != nil {
			log.Info(
				"failed to fetch calories burnt by id and date",
				zap.String("userId", userId.String()),
				zap.String("logDate", logDate),
				zap.Error(err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if *currValue+req.CaloriesBurnt < 0 {
			http.Error(w, "Resulting calories burnt can't be negative", http.StatusBadRequest)
			return
		}

		if err := lib.AddCaloriesBurntInTDEE(*userId, logDate, req.CaloriesBurnt); err != nil {
			log.Info(
				"failed to add burnt calories in tdee by id and date",
				zap.String("userId", userId.String()),
				zap.String("logDate", logDate),
				zap.Error(err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if req.CaloriesConsumed != 0.00 {
		currValue, err := lib.FetchCaloriesConsumedForTheDay(*userId, logDate)
		if err != nil {
			log.Info(
				"failed to fetch calories consumed by id and date",
				zap.String("userId", userId.String()),
				zap.String("logDate", logDate),
				zap.Error(err),
			)

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if *currValue+req.CaloriesConsumed < 0 {
			http.Error(w, "Resulting calories consumed can't be negative", http.StatusBadRequest)
			return
		}
	}

	if err := lib.UpdateCalorieLog(*userId, logDate, req.CaloriesConsumed, req.CaloriesBurnt); err != nil {
		log.Info(
			"failed to update calorie log by id and date",
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
