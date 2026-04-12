package handlers

import (
	"errors"

	"newco-go-reporting-service/internal/dto"
	"newco-go-reporting-service/internal/handlers/middleware"
	"newco-go-reporting-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type BranchDashboardHandler struct {
	service services.BranchDashboardService
}

func NewBranchDashboardHandler(service services.BranchDashboardService) *BranchDashboardHandler {
	return &BranchDashboardHandler{
		service: service,
	}
}

func (h *BranchDashboardHandler) Summary(c *fiber.Ctx) error {
	value := c.Locals(middleware.AccessContextKey)
	access, ok := value.(*dto.AccessContext)
	if !ok || access == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "access context missing",
		})
	}

	response, err := h.service.GetSummary(c.UserContext(), access)
	if err != nil {
		if errors.Is(err, services.ErrNoBranchScope) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"detail": "no branch scope available",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "failed to load branch dashboard summary",
		})
	}

	return c.JSON(response)
}

func (h *BranchDashboardHandler) BatchTrends(c *fiber.Ctx) error {
	value := c.Locals(middleware.AccessContextKey)
	access, ok := value.(*dto.AccessContext)
	if !ok || access == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "access context missing",
		})
	}

	items, err := h.service.GetBatchTrends(c.UserContext(), access)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "failed to load batch trends",
		})
	}

	return c.JSON(fiber.Map{
		"series": items,
	})
}

func (h *BranchDashboardHandler) RoleDistribution(c *fiber.Ctx) error {
	value := c.Locals(middleware.AccessContextKey)
	access, ok := value.(*dto.AccessContext)
	if !ok || access == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "access context missing",
		})
	}

	items, err := h.service.GetRoleDistribution(c.UserContext(), access)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "failed to load role distribution",
		})
	}

	return c.JSON(fiber.Map{
		"branchRoles": items,
	})
}

func (h *BranchDashboardHandler) RecentBatches(c *fiber.Ctx) error {
	value := c.Locals(middleware.AccessContextKey)
	access, ok := value.(*dto.AccessContext)
	if !ok || access == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"detail": "access context missing",
		})
	}

	items, err := h.service.GetRecentBatches(c.UserContext(), access)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": "failed to load recent batches",
		})
	}

	return c.JSON(fiber.Map{
		"items": items,
	})
}
