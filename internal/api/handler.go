package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"calometer/internal/db"

	"github.com/google/uuid"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
}

// Test endpoint
func GreetingHandler(w http.ResponseWriter, r *http.Request) {
	var greeting string
	if err := db.GetPool().QueryRow(context.Background(), "SELECT 'Hello, World'").Scan(&greeting); err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, greeting)
}

type User struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Gender string `json:"gender"`
}

func UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON body
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate the input
	if len(user.Gender) != 1 || user.Name == "" || user.Age <= 0 || user.Height <= 0 || user.Weight <= 0 {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	// Prepare the SQL query
	qStr := `
		INSERT INTO users (name, age, height, weight, gender)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var userId uuid.UUID
	if err := db.GetPool().QueryRow(
		context.Background(),
		qStr, user.Name,
		user.Age,
		user.Height,
		user.Weight,
		user.Gender,
	).Scan(&userId); err != nil {
		http.Error(w, "Failed to add to user", http.StatusInternalServerError)
		return
	}

	// Create the response with code and data
	response := Response{
		Code: http.StatusOK,
		Data: map[string]uuid.UUID{"user_id": userId},
	}

	// Set the response header and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
