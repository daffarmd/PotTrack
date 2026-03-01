package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pottrack/backend/internal/config"
	"github.com/pottrack/backend/internal/server"
)

func main() {
	// load .env if exists
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	engine := gin.Default()
	srv := server.New(engine, cfg)

	if err := srv.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
