package main

import (
	"log"

	"learn/internal/di"
)

func main() {
	if err := di.RunApp(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
