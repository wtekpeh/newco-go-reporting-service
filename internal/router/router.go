package router

import (
	"log"

	"newco-go-reporting-service/internal/config"
	"newco-go-reporting-service/internal/handlers"
	"newco-go-reporting-service/internal/handlers/middleware"
	"newco-go-reporting-service/internal/repositories"
	"newco-go-reporting-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(app *fiber.App, cfg *config.Config, pool *pgxpool.Pool) {
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
}
