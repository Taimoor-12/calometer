package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"calometer/internal/db"
	"calometer/internal/lib"

	"github.com/google/uuid"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" || user.Name == "" {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	doesExist, err := lib.DoesUserExists(user.Email)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if *doesExist {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	password, err := lib.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// Save user to the database
	qStr := `
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id
  `

	var userId uuid.UUID
	if err := db.GetPool().QueryRow(context.Background(), qStr, user.Name, user.Email, password).Scan(&userId); err != nil {
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
