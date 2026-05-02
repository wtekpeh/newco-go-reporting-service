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
	)

	// No empty sub-groups here.
	// Empty groups can behave like broad /reports middleware.

	reports.Get("/executive-summary", authMiddleware.RequireExecutive(), reportHandler.ExecutiveSummary)
	reports.Get("/recent-batches", authMiddleware.RequireExecutive(), reportHandler.RecentBatches)
	reports.Get("/batch-trends", authMiddleware.RequireExecutive(), reportHandler.BatchTrends)
	reports.Get("/branch-summary", authMiddleware.RequireExecutive(), reportHandler.BranchSummary)
	reports.Get("/role-distribution", authMiddleware.RequireExecutive(), reportHandler.RoleDistribution)
	reports.Get("/user-growth", authMiddleware.RequireExecutive(), reportHandler.UserGrowth)
	reports.Get("/staff-summary", authMiddleware.RequireExecutive(), reportHandler.StaffSummary)
	reports.Get("/batch-summary", authMiddleware.RequireExecutive(), reportHandler.BatchSummary)
	reports.Get("/branch-trends", authMiddleware.RequireExecutive(), reportHandler.BranchTrends)
	reports.Get("/branches", authMiddleware.RequireExecutive(), reportHandler.Branches)
	reports.Get("/ingredient-categories/daily", authMiddleware.RequireExecutive(), reportHandler.IngredientCategoryDaily)
	reports.Get("/ingredient-categories/daily/export/excel", authMiddleware.RequireExecutive(), reportHandler.ExportIngredientCategoryDailyExcel)

	reports.Get("/batches/:id/export/excel", authMiddleware.RequireBatchAccess(), reportHandler.ExportBatchDetailExcel)
	reports.Get("/batches/:id/export/pdf", authMiddleware.RequireBatchAccess(), reportHandler.ExportBatchDetailPDF)

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
