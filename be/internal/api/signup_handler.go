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
	resp := Response{}
	resp.Code = make(map[int]string)

	cookie, err := r.Cookie("token")
	if err == nil {
		// Validate the JWT
		if err := lib.ValidateToken(cookie.Value); err == nil {
			// Token is valid, return a success response
			w.WriteHeader(http.StatusOK)
			resp.Code[http.StatusOK] = "Logged in successfully."
			json.NewEncoder(w).Encode(&resp)
			return
		}
	}

	var user SignupHandlerReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Info("failed to decode incoming json")
		resp.Code[http.StatusBadRequest] = "Please enter correct details"
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if user.Username == "" || user.Password == "" || user.Name == "" {
		log.Info("invalid input data")
		resp.Code[http.StatusBadRequest] = "Please enter correct details."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	doesExist, err := lib.DoesUserExists(user.Username)
	if err != nil {
		log.Info(
			"failed to check user's existence by username",
			zap.String("username", user.Username),
		)
		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if *doesExist {
		log.Info("username already exists")
		resp.Code[http.StatusConflict] = "Username already exists"
		json.NewEncoder(w).Encode(&resp)
		return
	}

	password, err := lib.HashPassword(user.Password)
	if err != nil {
		log.Info(
			"failed to hash user's password",
			zap.String("username", user.Username),
		)
		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
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
			resp.Code[http.StatusConflict] = "Username already exists."
			json.NewEncoder(w).Encode(&resp)
			return
		}

		log.Info(
			"failed to create the user",
			zap.String("username", user.Username),
		)
		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	w.WriteHeader(http.StatusOK)

	data := &SignupHandlerResp{
		UserId: userId,
	}

	resp.Code[http.StatusOK] = "OK"
	resp.Data = data
	json.NewEncoder(w).Encode(&resp)
}
