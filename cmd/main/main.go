package main

import (
	"log/slog"
	"os"

	"learn/internal/di"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := di.RunApp(logger); err != nil {
		logger.Error("failed to run app", slog.Any("error", err))
		os.Exit(1)
	}
}
