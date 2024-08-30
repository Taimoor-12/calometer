package api

import (
	"context"
	"fmt"
	"net/http"

	"calometer/internal/db"
)

func GreetingHandler(w http.ResponseWriter, r *http.Request) {
	var greeting string
	if err := db.GetPool().QueryRow(context.Background(), "SELECT 'Hello, World'").Scan(&greeting); err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, greeting)
}
