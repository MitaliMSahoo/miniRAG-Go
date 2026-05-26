package main

import (
	"log"
	"miniRAGServer-Go/config"
	miniragserver "miniRAGServer-Go/miniRAGServer"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	log.Printf("Starting server...")
	log.Printf("Loading environment variables...")
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: failed to load .env file: %w", err)
	}
	cfg := config.Load()
	s := miniragserver.ServerInit()
	log.Printf("Server starting on %s", cfg.App.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.App.Port, s.Handler()))
}
