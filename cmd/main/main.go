package main

import (
	"log/slog"
	"os"

	_ "learn/docs"
	"learn/internal/di"
)

// @title Subscription Service API
// @version 1.0
// @description REST API for managing user subscriptions
// @host localhost:8090
// @BasePath /api/v1
// @schemes http
// @accept json
// @produce json
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := di.RunApp(logger); err != nil {
		logger.Error("failed to run app", slog.Any("error", err))
		os.Exit(1)
	}
}
