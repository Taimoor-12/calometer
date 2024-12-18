package api

import (
	"calometer/internal/lib"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

type LoginHandlerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginHandlerResp struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	var req LoginHandlerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Info("failed to decode incoming json")
		resp.Code[http.StatusBadRequest] = "Invalid JSON data."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if req.Username == "" && req.Password == "" {
		log.Info("invalid input data")
		resp.Code[http.StatusBadRequest] = "Please enter correct details."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	userId, err := lib.GetUserIdByUsername(req.Username)
	if err == nil && userId == nil {
		resp.Code[http.StatusUnauthorized] = "Username or password is incorrect."
		json.NewEncoder(w).Encode(&resp)
		return
	} else if err != nil {
		log.Info(
			"failed to get user id by username",
			zap.String("username", req.Username),
		)

		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	exists, err := lib.DoesUserExists(req.Username)
	if err != nil {
		log.Info(
			"failed to check user's existence by username",
			zap.String("username", req.Username),
		)
		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if !*exists {
		resp.Code[http.StatusUnauthorized] = "Username or password is incorrect."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	passwordHash, err := lib.GetHashedPass(req.Username)
	if err != nil {
		log.Info(
			"failed to fetch user's hashed password",
			zap.String("username", req.Username),
		)
		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	if err := lib.CheckPasswordValidity(req.Password, passwordHash); err != nil {
		log.Info(
			"failed to check password validity",
			zap.String("username", req.Username),
		)
		resp.Code[http.StatusUnauthorized] = "Username or password is incorrect."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	token, err := lib.GenerateJWT(*userId, req.Username)
	if err != nil {
		log.Info(
			"failed to generate JWT for user id",
			zap.String("userId", userId.String()),
			zap.String("username", req.Username),
		)
		resp.Code[http.StatusInternalServerError] = "Something went wrong, please try again."
		json.NewEncoder(w).Encode(&resp)
		return
	}

	// Set the JWT as an HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		Expires:  time.Now().Add(time.Hour * 24),
	})

	w.WriteHeader(http.StatusOK)
	resp.Code[http.StatusOK] = "Logged in successfully."
	json.NewEncoder(w).Encode(&resp)
}
