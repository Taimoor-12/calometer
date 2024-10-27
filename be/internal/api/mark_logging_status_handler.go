package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type MarkLoggingStatusReq struct {
	Status  string    `json:"status"`
	LogDate time.Time `json:"log_date"`
}

func MarkLoggingStatusHandler(w http.ResponseWriter, r *http.Request) {
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

	var req MarkLoggingStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		resp.Code[http.StatusBadRequest] = "Invalid JSON."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	logDate := req.LogDate.Format("2006-01-02")

	if err := lib.MarkLoggingStatus(*userId, logDate, req.Status); err != nil {
		log.Info(
			"failed to mark logging status by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	calorieLogId, err := lib.GetCalorieLogId(*userId, logDate)
	if err != nil {
		log.Info(
			"failed to get calorie log by id and date",
			zap.String("userId", userId.String()),
			zap.String("logDate", logDate),
			zap.Error(err),
		)

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if req.Status == "D" {
		caloricBalance, err := lib.CalculateCaloricBalanceForTheDay(*userId, logDate)
		if err != nil {
			log.Info(
				"failed to calculate caloric balance by id and date",
				zap.String("userId", userId.String()),
				zap.String("logDate", logDate),
				zap.Error(err),
			)

			resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
			json.NewEncoder(w).Encode(&resp)
			return
		}

		if err := lib.AddCaloricBalanceForTheDay(*calorieLogId, *caloricBalance); err != nil {
			log.Info(
				"failed to add caloric balance by log id",
				zap.String("logId", calorieLogId.String()),
				zap.Error(err),
			)

			resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
			json.NewEncoder(w).Encode(&resp)
			return
		}
	} else {
		if err := lib.ResetCaloricBalanceForTheDay(*calorieLogId); err != nil {
			log.Info(
				"failed to reset caloric balance by log id",
				zap.String("logId", calorieLogId.String()),
				zap.Error(err),
			)

			resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
			json.NewEncoder(w).Encode(&resp)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	resp.Code[http.StatusOK] = "OK"
	json.NewEncoder(w).Encode(&resp)
}
