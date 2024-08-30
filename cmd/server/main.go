package main

import (
	"calometer/internal/api"
	"calometer/internal/db"
	"fmt"
	"log"
	"net/http"
)

func main() {
	pool, err := db.Init()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Close(pool)

	router := api.SetupRouter()

	fmt.Println("Starting server on: :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
