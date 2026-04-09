package handlers

import (
	"newco-go-reporting-service/internal/dto"
	"newco-go-reporting-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	Service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{
		Service: service,
	}
}

func (h *ReportHandler) ExecutiveSummary(c *fiber.Ctx) error {
	totalUsers, err := h.Service.TotalUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch total users",
			Error:   err.Error(),
		})
	}

	totalActiveUsers, err := h.Service.TotalActiveUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch total active users",
			Error:   err.Error(),
		})
	}

	response := dto.ExecutiveSummaryResponse{
		Message:          "executive summary endpoint ready",
		TotalUsers:       totalUsers,
		TotalActiveUsers: totalActiveUsers,
	}

	return c.JSON(response)
}
