package app

import (
	"context"
	"time"

	notificationsrepo "newco-go-reporting-service/internal/notifications/repositories"
	notificationsservice "newco-go-reporting-service/internal/notifications/services"
	notificationsworker "newco-go-reporting-service/internal/notifications/worker"

	"newco-go-reporting-service/internal/config"
	"newco-go-reporting-service/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg *config.Config, pool *pgxpool.Pool) *fiber.App {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://vite.williamtekpeh.com,http://localhost:5173,http://127.0.0.1:5173",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	router.Register(app, cfg, pool)

	eventRepo := notificationsrepo.NewEventRepository(pool)
	notificationRepo := notificationsrepo.NewNotificationRepository(pool)
	eventProcessor := notificationsservice.NewEventProcessorService(
		eventRepo,
		notificationRepo,
	)
	eventWorker := notificationsworker.NewEventWorker(
		eventProcessor,
		5*time.Second,
		20,
	)

	go eventWorker.Start(context.Background())

	return app
}
