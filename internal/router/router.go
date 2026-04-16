package router

import (
	"log"

	"newco-go-reporting-service/internal/config"
	"newco-go-reporting-service/internal/handlers"
	"newco-go-reporting-service/internal/handlers/middleware"
	"newco-go-reporting-service/internal/repositories"
	"newco-go-reporting-service/internal/services"

	notificationhandlers "newco-go-reporting-service/internal/notifications/handlers"
	notificationsrealtime "newco-go-reporting-service/internal/notifications/realtime"
	notificationsrepo "newco-go-reporting-service/internal/notifications/repositories"
	notificationsservice "newco-go-reporting-service/internal/notifications/services"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(
	app *fiber.App,
	cfg *config.Config,
	pool *pgxpool.Pool,
	hub *notificationsrealtime.Hub,
) {
	_ = hub
	healthHandler := handlers.NewHealthHandler(pool)

	reportRepo := repositories.NewReportRepository(pool)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	accessRepo := repositories.NewAccessRepository(pool)
	branchAccessRepo := repositories.NewBranchAccessRepository(pool)
	accessService := services.NewAccessService(accessRepo, branchAccessRepo)

	branchDashboardRepo := repositories.NewBranchDashboardRepository(pool)
	branchDashboardService := services.NewBranchDashboardService(branchDashboardRepo)
	branchDashboardHandler := handlers.NewBranchDashboardHandler(branchDashboardService)

	notificationRepo := notificationsrepo.NewNotificationRepository(pool)
	notificationService := notificationsservice.NewNotificationService(notificationRepo)
	notificationHandler := notificationhandlers.NewNotificationHandler(notificationService)

	notificationWSHandler := notificationhandlers.NewNotificationWebSocketHandler(hub)

	authMiddleware, err := middleware.NewAuthMiddleware(cfg, accessService)
	if err != nil {
		log.Fatal(err)
	}

	app.Get("/", healthHandler.Check)
	app.Get("/health", healthHandler.Check)

	reports := app.Group("/reports",
		authMiddleware.RequireAuth(),
		authMiddleware.RequireExecutive(),
	)

	reports.Get("/executive-summary", reportHandler.ExecutiveSummary)
	reports.Get("/recent-batches", reportHandler.RecentBatches)
	reports.Get("/batch-trends", reportHandler.BatchTrends)
	reports.Get("/branch-summary", reportHandler.BranchSummary)
	reports.Get("/role-distribution", reportHandler.RoleDistribution)
	reports.Get("/user-growth", reportHandler.UserGrowth)
	reports.Get("/staff-summary", reportHandler.StaffSummary)
	reports.Get("/batch-summary", reportHandler.BatchSummary)
	reports.Get("/branch-trends", reportHandler.BranchTrends)
	reports.Get("/branches", reportHandler.Branches)

	branchDashboard := app.Group(
		"/branch-dashboard",
		authMiddleware.RequireAuth(),
		authMiddleware.RequireBranchManager(),
	)

	branchDashboard.Get("/summary", branchDashboardHandler.Summary)
	branchDashboard.Get("/batch-trends", branchDashboardHandler.BatchTrends)
	branchDashboard.Get("/role-distribution", branchDashboardHandler.RoleDistribution)
	branchDashboard.Get("/recent-batches", branchDashboardHandler.RecentBatches)

	notifications := app.Group(
		"/notifications",
		authMiddleware.RequireAuth(),
	)

	notifications.Get("/", notificationHandler.ListMyNotifications)
	notifications.Get("/unread-count", notificationHandler.GetMyUnreadCount)
	notifications.Patch("/:id/read", notificationHandler.MarkMyNotificationAsRead)
	notifications.Patch("/read-all", notificationHandler.MarkAllMyNotificationsAsRead)

	app.Use(
		"/ws/notifications",
		authMiddleware.RequireWebSocketAuth(),
	)

	app.Get(
		"/ws/notifications",
		websocket.New(notificationWSHandler.Handle),
	)
}
