package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type AddBodyDetailsReq struct {
	Age       int     `json:"age"`
	Weight_kg float64 `json:"weight"`
	Height_cm int     `json:"height"`
	Gender    string  `json:"gender"`
	Goal      string  `json:"goal,omitempty"`
}

func AddBodyDetailsHandler(w http.ResponseWriter, r *http.Request) {
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

	var req AddBodyDetailsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		resp.Code[http.StatusBadRequest] = "Invalid JSON data."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	userId, err := lib.ExtractUserIdFromToken(tokenStr)
	if err != nil {
		log.Info(
			"failed to extract username from token",
			zap.String("tokenStr", tokenStr),
			zap.Error(err),
		)

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if err := lib.AddUserBodyDetails(
		*userId,
		req.Age,
		req.Height_cm,
		req.Weight_kg,
		req.Gender,
	); err != nil {
		log.Info(
			"failed to add user body details by id",
			zap.String("userId", userId.String()),
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
