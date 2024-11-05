package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type DoBodyDetailsExistResp struct {
	Exists bool `json:"exists"`
}

func DoBodyDetailsExistHandler(w http.ResponseWriter, r *http.Request) {
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

	exists, err := lib.DoesBodyDetailsExist(*userId)
	if err != nil {
		log.Info(
			"failed to check existence of body details by user id",
			zap.String("userId", userId.String()),
			zap.Error(err),
		)

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp.Code[http.StatusOK] = "OK"
	data := &DoBodyDetailsExistResp{
		Exists: *exists,
	}
	resp.Data = data
	json.NewEncoder(w).Encode(&resp)
}
