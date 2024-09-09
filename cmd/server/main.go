package main

import (
	"calometer/internal/api"
	"calometer/internal/db"
	"calometer/internal/logger"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	log := logger.GetLogger()

	log.Info(
		"Application Started",
		zap.String("name", "calometer"),
	)

	// Init database
	pool, err := db.Init()
	if err != nil {
		log.Fatal("Failed to initliaze database", zap.Error(err))
	}
	defer db.Close(pool)

	// Init server
	router := api.SetupRouter()

	fmt.Println("Starting server on: :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed to start", zap.Error(err))
	}
}
