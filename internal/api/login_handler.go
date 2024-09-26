package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Check for JWT in the request
	tokenStr := lib.ExtractTokenFromHeader(r)
	if tokenStr != "" {
		// Step 2: Validate the JWT
		if err := lib.ValidateToken(tokenStr); err == nil {
			// Token is valid, return a success response
			w.WriteHeader(http.StatusOK)
			resp := Response{
				Code: http.StatusOK,
				Data: map[string]string{"token": tokenStr},
			}
			json.NewEncoder(w).Encode(&resp)
			return
		}
	}

	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info("failed to decode incoming json")
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if req.Username == "" && req.Password == "" {
		log.Info("invalid input data")
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	userId, err := lib.GetUserIdByUsername(req.Username)
	if err != nil {
		log.Fatal(
			"failed to get user id by username",
			zap.String("username", req.Username),
		)
	}

	exists, err := lib.DoesUserExists(req.Username)
	if err != nil {
		log.Info(
			"failed to check user's existence by username",
			zap.String("username", req.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !*exists {
		http.Error(w, "Username or Password is incorrect", http.StatusUnauthorized)
		return
	}

	passwordHash, err := lib.GetHashedPass(req.Username)
	if err != nil {
		log.Info(
			"failed to fetch user's hashed password",
			zap.String("username", req.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := lib.CheckPasswordValidity(req.Password, passwordHash); err != nil {
		log.Info(
			"failed to check password validity",
			zap.String("username", req.Username),
		)
		http.Error(w, "Username or Password is incorrect", http.StatusUnauthorized)
		return
	}

	token, err := lib.GenerateJWT(*userId, req.Username)
	if err != nil {
		log.Info(
			"failed to generate JWT for user id",
			zap.String("userId", userId.String()),
			zap.String("username", req.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{
		Code: http.StatusOK,
		Data: map[string]string{"token": token},
	}
	json.NewEncoder(w).Encode(&resp)
}
