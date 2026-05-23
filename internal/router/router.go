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

	aihandlers "newco-go-reporting-service/internal/ai/handlers"
	aiservices "newco-go-reporting-service/internal/ai/services"
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

	ollamaService := aiservices.NewOllamaService(
		cfg.OllamaBaseURL,
		cfg.OllamaModel,
	)

	executiveContextBuilder := aiservices.NewExecutiveContextBuilder(
		reportService,
	)

	chatMemoryStore := aiservices.NewChatMemoryStore()

	intentClassifier := aiservices.NewIntentClassifier(
		ollamaService,
		chatMemoryStore,
	)

	aiChatHandler := aihandlers.NewAIChatHandler(
		intentClassifier,
		ollamaService,
		reportService,
		chatMemoryStore,
		executiveContextBuilder,
	)

	executiveAIHandler := aihandlers.NewExecutiveAIHandler(
		ollamaService,
		executiveContextBuilder,
	)

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
	reports.Get("/top-recipe-variance", authMiddleware.RequireExecutive(), reportHandler.GetTopRecipeVariance)
	reports.Get("/branch-trends", authMiddleware.RequireExecutive(), reportHandler.BranchTrends)
	reports.Get("/daily-plan-trends", authMiddleware.RequireExecutive(), reportHandler.DailyPlanTrends)
	reports.Get("/recent-daily-plans", authMiddleware.RequireExecutive(), reportHandler.RecentDailyPlans)
	reports.Get("/branches", authMiddleware.RequireExecutive(), reportHandler.Branches)
	reports.Get("/ingredient-categories/daily", authMiddleware.RequireExecutive(), reportHandler.IngredientCategoryDaily)
	reports.Get("/ingredient-categories/daily/export/excel", authMiddleware.RequireExecutive(), reportHandler.ExportIngredientCategoryDailyExcel)
	reports.Get(
		"/ingredient-categories/daily/export/pdf", authMiddleware.RequireExecutive(),
		reportHandler.ExportIngredientCategoryDailyPDF,
	)

	reports.Get("/batches/:id/export/excel", authMiddleware.RequireBatchAccess(), reportHandler.ExportBatchDetailExcel)
	reports.Get("/batches/:id/export/pdf", authMiddleware.RequireBatchAccess(), reportHandler.ExportBatchDetailPDF)
	reports.Get("/daily-plans/:id/export/pdf", authMiddleware.RequireBatchAccess(), reportHandler.DailyPlanRequisitionPDF)

	branchDashboard := app.Group(
		"/branch-dashboard",
		authMiddleware.RequireAuth(),
		authMiddleware.RequireBranchManager(),
	)

	branchDashboard.Get("/summary", branchDashboardHandler.Summary)
	branchDashboard.Get("/batch-trends", branchDashboardHandler.BatchTrends)
	branchDashboard.Get("/daily-plan-trends", branchDashboardHandler.DailyPlanTrends)
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

	ai := app.Group(
		"/ai",
		authMiddleware.RequireAuth(),
		authMiddleware.RequireExecutive(),
	)

	ai.Post(
		"/executive-summary",
		executiveAIHandler.GenerateSummary,
	)

	ai.Post(
		"/chat",
		aiChatHandler.Chat,
	)

	app.Use(
		"/ws/notifications",
		authMiddleware.RequireWebSocketAuth(),
	)

	app.Get(
		"/ws/notifications",
		websocket.New(notificationWSHandler.Handle),
	)
}
