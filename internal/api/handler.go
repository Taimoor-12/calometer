package api

import (
	"context"
	"encoding/json"
	"net/http"

	"calometer/internal/db"
	"calometer/internal/lib"
	"calometer/internal/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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
		log.Info("username already exists")
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
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			log.Info(
				"username already exists",
				zap.String("username", user.Username),
			)
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}

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

type AddBodyDetailsReq struct {
	Age       int     `json:"age"`
	Weight_kg float64 `json:"weight"`
	Height_cm int     `json:"height"`
	Gender    string  `json:"gender"`
}

func AddBodyDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Check for JWT in the request
	tokenStr := lib.ExtractTokenFromHeader(r)
	if tokenStr == "" {
		// No token found in the request header
		http.Error(w, "User not authenticated. Please provide a valid token.", http.StatusUnauthorized)
		return
	}

	// Validate the JWT
	if err := lib.ValidateToken(tokenStr); err != nil {
		// Token is invalid
		http.Error(w, "Invalid token. Please authenticate again.", http.StatusUnauthorized)
		return
	}

	var req AddBodyDetailsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info(
			"failed to decode incoming json",
			zap.Error(err),
		)

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userId, err := lib.ExtractUserIdFromToken(tokenStr)
	if err != nil {
		log.Info(
			"failed to extract username from token",
			zap.String("tokenStr", tokenStr),
			zap.Error(err),
		)
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

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{
		Code: http.StatusOK,
	}
	json.NewEncoder(w).Encode(&resp)
}
