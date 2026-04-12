package main

import (
	"log"

	"github.com/joho/godotenv"

	"newco-go-reporting-service/internal/app"
	"newco-go-reporting-service/internal/config"
	"newco-go-reporting-service/internal/db"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	cfg := config.Load()

	pool, err := db.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	server := app.New(cfg, pool)

	log.Fatal(server.Listen(":" + cfg.AppPort))
}
