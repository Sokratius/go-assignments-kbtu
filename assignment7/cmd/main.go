package main

import (
	"log"
	"practice-7/config"
	"practice-7/internal/app"
)

func main() {
	cfg := config.New()

	a, err := app.New(cfg)
	if err != nil {
		log.Fatalf("app init error: %v", err)
	}

	if err := a.Run(); err != nil {
		log.Fatalf("app run error: %v", err)
	}
}
