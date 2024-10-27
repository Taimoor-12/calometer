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

	var req DeleteCalorieLogReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		resp.Code[http.StatusBadRequest] = "Invalid JSON data."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	var logDate string
	// Check if the logDate is in the future
	if req.LogDate.After(time.Now()) {
		resp.Code[http.StatusBadRequest] = "Log date cannot be a future date."
		json.NewEncoder(w).Encode(&resp)
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

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if !*exists {
		resp.Code[http.StatusConflict] = "No log exists for this day."
		json.NewEncoder(w).Encode(&resp)
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

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if err := lib.DeleteCaloricBalanceByLogId(*logId); err != nil {
		log.Info(
			"failed to delete caloric balance by log id",
			zap.String("logId", logId.String()),
			zap.Error(err),
		)

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if err := lib.DeleteCalorieLog(*userId, logDate); err != nil {
		log.Info(
			"failed to delete log by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp.Code[http.StatusOK] = "OK"
	json.NewEncoder(w).Encode(&resp)
}
