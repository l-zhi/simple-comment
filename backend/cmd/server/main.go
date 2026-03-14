package main

import (
	"log"
	"os"

	"simple-comment/internal/config"
	"simple-comment/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("failed to load config: %v", err)
		os.Exit(1)
	}

	r := router.New(cfg)

	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Printf("server exited with error: %v", err)
		os.Exit(1)
	}
}