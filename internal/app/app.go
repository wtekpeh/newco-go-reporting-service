package app

import (
	"newco-go-reporting-service/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pool *pgxpool.Pool) *fiber.App {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())

	router.Register(app, pool)

	return app
}
