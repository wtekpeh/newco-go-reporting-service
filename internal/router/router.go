package router

import (
	"newco-go-reporting-service/internal/handlers"
	"newco-go-reporting-service/internal/repositories"
	"newco-go-reporting-service/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(app *fiber.App, pool *pgxpool.Pool) {
	healthHandler := handlers.NewHealthHandler(pool)

	reportRepo := repositories.NewReportRepository(pool)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	app.Get("/", healthHandler.Check)
	app.Get("/health", healthHandler.Check)

	reports := app.Group("/reports")
	reports.Get("/executive-summary", reportHandler.ExecutiveSummary)
	reports.Get("/recent-batches", reportHandler.RecentBatches)
	reports.Get("/batch-trends", reportHandler.BatchTrends)
	reports.Get("/branch-summary", reportHandler.BranchSummary)
	reports.Get("/role-distribution", reportHandler.RoleDistribution)
	reports.Get("/user-growth", reportHandler.UserGrowth)
	reports.Get("/staff-summary", reportHandler.StaffSummary)
	reports.Get("/batch-summary", reportHandler.BatchSummary)
	reports.Get("/branch-trends", reportHandler.BranchTrends)
}
