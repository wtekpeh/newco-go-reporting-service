package handlers

import (
	"newco-go-reporting-service/internal/dto"
	"newco-go-reporting-service/internal/services"
	"strconv"
	"time"

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
	filters := h.parseFilters(c)
	totalUsers, err := h.Service.TotalUsers(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch total users",
			Error:   err.Error(),
		})
	}

	totalActiveUsers, err := h.Service.TotalActiveUsers(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch total active users",
			Error:   err.Error(),
		})
	}

	totalBranches, err := h.Service.TotalBranches()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch total branches",
			Error:   err.Error(),
		})
	}

	totalBatches, err := h.Service.TotalBatches(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch total batches",
			Error:   err.Error(),
		})
	}

	batchesThisWeek, err := h.Service.BatchesThisWeek(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch batches this week",
			Error:   err.Error(),
		})
	}

	batchesThisMonth, err := h.Service.BatchesThisMonth(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch batches this month",
			Error:   err.Error(),
		})
	}

	mostActiveBranch, err := h.Service.MostActiveBranch(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch most active branch",
			Error:   err.Error(),
		})
	}

	largestBranch, err := h.Service.LargestBranch(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch largest branch",
			Error:   err.Error(),
		})
	}

	mostUsedRecipe, err := h.Service.MostUsedRecipe(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch most used recipe",
			Error:   err.Error(),
		})
	}

	averageBatchesPerBranch, err := h.Service.AverageBatchesPerBranch(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch average batches per branch",
			Error:   err.Error(),
		})
	}

	peakBatchDay, err := h.Service.PeakBatchDay(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch peak batch day",
			Error:   err.Error(),
		})
	}

	response := dto.ExecutiveSummaryResponse{
		Message: "executive summary endpoint ready",
		Kpis: dto.ExecutiveKpisDTO{
			TotalUsers:       totalUsers,
			TotalActiveUsers: totalActiveUsers,
			TotalBranches:    totalBranches,
			TotalBatches:     totalBatches,
			BatchesThisWeek:  batchesThisWeek,
			BatchesThisMonth: batchesThisMonth,
		},
		Highlights: dto.ExecutiveHighlightsDTO{
			MostActiveBranch:        mostActiveBranch,
			LargestBranch:           largestBranch,
			MostUsedRecipe:          mostUsedRecipe,
			AverageBatchesPerBranch: averageBatchesPerBranch,
			PeakBatchDay:            peakBatchDay,
		},
	}

	return c.JSON(response)
}

func (h *ReportHandler) RecentBatches(c *fiber.Ctx) error {
	filters := h.parseFilters(c)
	items, err := h.Service.RecentBatches(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch recent batches",
			Error:   err.Error(),
		})
	}

	response := dto.RecentBatchesResponse{
		Message: "recent batches fetched successfully",
		Items:   items,
	}

	return c.JSON(response)
}

func (h *ReportHandler) BatchTrends(c *fiber.Ctx) error {
	filters := h.parseFilters(c)
	series, err := h.Service.BatchTrends(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch batch trends",
			Error:   err.Error(),
		})
	}

	response := dto.BatchTrendsResponse{
		Message: "batch trends fetched successfully",
		Series:  series,
	}

	return c.JSON(response)
}

func (h *ReportHandler) BranchSummary(c *fiber.Ctx) error {
	filters := h.parseFilters(c)
	items, err := h.Service.BranchSummary(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch branch summary",
			Error:   err.Error(),
		})
	}

	response := dto.BranchSummaryResponse{
		Message: "branch summary fetched successfully",
		Items:   items,
	}

	return c.JSON(response)
}

func (h *ReportHandler) RoleDistribution(c *fiber.Ctx) error {
	items, err := h.Service.RoleDistribution()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch role distribution",
			Error:   err.Error(),
		})
	}

	response := dto.RoleDistributionResponse{
		Message: "role distribution fetched successfully",
		Items:   items,
	}

	return c.JSON(response)
}

func (h *ReportHandler) UserGrowth(c *fiber.Ctx) error {
	filters := h.parseFilters(c)
	series, err := h.Service.UserGrowth(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch user growth",
			Error:   err.Error(),
		})
	}

	response := dto.UserGrowthResponse{
		Message: "user growth fetched successfully",
		Series:  series,
	}

	return c.JSON(response)
}

func (h *ReportHandler) parseFilters(c *fiber.Ctx) dto.ReportFiltersDTO {
	return dto.ReportFiltersDTO{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		BranchID:  c.Query("branch_id"),
		GroupBy:   c.Query("group_by"),
	}
}

func (h *ReportHandler) StaffSummary(c *fiber.Ctx) error {
	filters := h.parseFilters(c)

	response, err := h.Service.StaffSummary(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch staff summary",
			Error:   err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *ReportHandler) BatchSummary(c *fiber.Ctx) error {
	filters := h.parseFilters(c)

	response, err := h.Service.BatchSummary(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch batch summary",
			Error:   err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *ReportHandler) BranchTrends(c *fiber.Ctx) error {
	filters := h.parseFilters(c)

	series, err := h.Service.BranchTrends(filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch branch trends",
			Error:   err.Error(),
		})
	}

	response := dto.BranchTrendsResponse{
		Message: "branch trends fetched successfully",
		Series:  series,
	}

	return c.JSON(response)
}

func (h *ReportHandler) Branches(c *fiber.Ctx) error {
	items, err := h.Service.GetBranches()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch branches",
			Error:   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"items": items,
	})
}

func (h *ReportHandler) IngredientCategoryDaily(c *fiber.Ctx) error {

	dateStr := c.Query("date")
	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "date query parameter is required",
			Error:   "use format YYYY-MM-DD",
		})
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "invalid date format",
			Error:   "use format YYYY-MM-DD",
		})
	}

	items, err := h.Service.GetIngredientCategoryDaily(c.Context(), date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch ingredient category daily report",
			Error:   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "ingredient category daily report fetched successfully",
		"items":   items,
	})
}

func (h *ReportHandler) ExportIngredientCategoryDailyExcel(c *fiber.Ctx) error {
	dateStr := c.Query("date")
	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "date query parameter is required",
			Error:   "use format YYYY-MM-DD",
		})
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "invalid date format",
			Error:   "use format YYYY-MM-DD",
		})
	}

	fileBytes, err := h.Service.ExportIngredientCategoryDailyExcel(c.Context(), date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to export ingredient category daily excel",
			Error:   err.Error(),
		})
	}

	filename := "ingredient_category_daily_" + dateStr + ".xlsx"

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.Send(fileBytes)
}

func (h *ReportHandler) ExportBatchDetailExcel(c *fiber.Ctx) error {
	batchID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "invalid batch id",
			Error:   "batch id must be a number",
		})
	}

	fileBytes, err := h.Service.ExportBatchDetailExcel(c.Context(), batchID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to export batch detail excel",
			Error:   err.Error(),
		})
	}

	filename := "batch_detail_" + strconv.FormatInt(batchID, 10) + ".xlsx"

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.Send(fileBytes)
}

func (h *ReportHandler) ExportBatchDetailPDF(c *fiber.Ctx) error {
	batchID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "invalid batch id",
			Error:   "batch id must be a number",
		})
	}

	fileBytes, err := h.Service.ExportBatchDetailPDF(c.Context(), batchID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to export batch detail pdf",
			Error:   err.Error(),
		})
	}

	filename := "batch_detail_" + strconv.FormatInt(batchID, 10) + ".pdf"

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	return c.Send(fileBytes)
}

func (h *ReportHandler) GetTopRecipeVariance(c *fiber.Ctx) error {
	startDate := c.Query("start_date", "")
	endDate := c.Query("end_date", "")
	branchID := c.Query("branch_id", "")

	result, err := h.Service.GetTopRecipeVariance(
		c.Context(),
		startDate,
		endDate,
		branchID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	return c.JSON(result)
}
