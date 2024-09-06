package api

import (
	"context"
	"encoding/json"
	"net/http"

	"calometer/internal/db"
	"calometer/internal/lib"
	"calometer/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var log = logger.GetLogger()

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

type SignupReq struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user SignupReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Info("failed to decode incoming json")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" || user.Name == "" {
		log.Info("invalid input data")
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	doesExist, err := lib.DoesUserExists(user.Username)
	if err != nil {
		log.Info(
			"failed to check user's existence by username",
			zap.String("username", user.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if *doesExist {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	password, err := lib.HashPassword(user.Password)
	if err != nil {
		log.Info(
			"failed to hash user's password",
			zap.String("username", user.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Save user to the database
	qStr := `
	INSERT INTO users (name, username, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id
  `

	var userId uuid.UUID
	if err := db.GetPool().QueryRow(context.Background(), qStr, user.Name, user.Username, password).Scan(&userId); err != nil {
		log.Info(
			"failed to create the user",
			zap.String("username", user.Username),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	response := Response{
		Code: http.StatusOK,
		Data: map[string]uuid.UUID{"user_id": userId},
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}

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

	token, err := lib.GenerateJWT(req.Username)
	if err != nil {
		log.Info(
			"failed to generate JWT",
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
