package api

import (
	"calometer/internal/db"
	"calometer/internal/lib"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type SignupHandlerReq struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupHandlerResp struct {
	UserId uuid.UUID `json:"u_id"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user SignupHandlerReq
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

	w.WriteHeader(http.StatusOK)

	data := &SignupHandlerResp{
		UserId: userId,
	}

	resp := Response{
		Code: http.StatusOK,
		Data: data,
	}
	json.NewEncoder(w).Encode(&resp)
}
