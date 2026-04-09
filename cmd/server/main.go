package main

import (
	"log"

	"newco-go-reporting-service/internal/app"
	"newco-go-reporting-service/internal/config"
	"newco-go-reporting-service/internal/db"
)

func main() {
	cfg := config.Load()

	pool, err := db.NewPostgresPool(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	server := app.New(pool)

	log.Fatal(server.Listen(":" + cfg.AppPort))
}
