package services

import (
	"bytes"
	"context"
	"newco-go-reporting-service/internal/dto"
	"newco-go-reporting-service/internal/repositories"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"

	"github.com/xuri/excelize/v2"
)

type ReportService struct {
	Repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{
		Repo: repo,
	}
}

func exportDisplayUnit(unit string) string {
	switch unit {
	case "g":
		return "kg"
	case "ml":
		return "L"
	case "pc":
		return "pcs"
	case "pcs":
		return "pcs"
	default:
		return unit
	}
}

func exportDisplayValue(unit string, value float64) float64 {
	switch unit {
	case "g", "ml":
		return value / 1000
	default:
		return value
	}
}

func exportFloat3(value float64) string {
	return strconv.FormatFloat(value, 'f', 3, 64)
}

func (s *ReportService) TotalUsers(filters dto.ReportFiltersDTO) (int, error) {
	return s.Repo.TotalUsers(filters)
}

func (s *ReportService) TotalActiveUsers(filters dto.ReportFiltersDTO) (int, error) {
	return s.Repo.TotalActiveUsers(filters)
}

func (s *ReportService) TotalBranches() (int, error) {
	return s.Repo.TotalBranches()
}

func (s *ReportService) TotalBatches(filters dto.ReportFiltersDTO) (int, error) {
	return s.Repo.TotalBatches(filters)
}

func (s *ReportService) BatchesThisWeek(filters dto.ReportFiltersDTO) (int, error) {
	return s.Repo.BatchesThisWeek(filters)
}

func (s *ReportService) BatchesThisMonth(filters dto.ReportFiltersDTO) (int, error) {
	return s.Repo.BatchesThisMonth(filters)
}

func (s *ReportService) MostActiveBranch(filters dto.ReportFiltersDTO) (*dto.BranchActivityDTO, error) {
	return s.Repo.MostActiveBranch(filters)
}

func (s *ReportService) LargestBranch(filters dto.ReportFiltersDTO) (*dto.BranchStaffDTO, error) {
	return s.Repo.LargestBranch(filters)
}

func (s *ReportService) MostUsedRecipe(filters dto.ReportFiltersDTO) (*dto.RecipeUsageDTO, error) {
	return s.Repo.MostUsedRecipe(filters)
}

func (s *ReportService) AverageBatchesPerBranch(filters dto.ReportFiltersDTO) (*dto.AverageMetricDTO, error) {
	return s.Repo.AverageBatchesPerBranch(filters)
}

func (s *ReportService) PeakBatchDay(filters dto.ReportFiltersDTO) (*dto.PeakBatchDayDTO, error) {
	return s.Repo.PeakBatchDay(filters)
}

func (s *ReportService) RecentBatches(filters dto.ReportFiltersDTO) ([]dto.RecentBatchItemDTO, error) {
	return s.Repo.RecentBatches(filters)
}

func (s *ReportService) BatchTrends(filters dto.ReportFiltersDTO) ([]dto.BatchTrendPointDTO, error) {
	if filters.GroupBy == "" {
		filters.GroupBy = "day"
	}

	return s.Repo.BatchTrends(filters)
}

func (s *ReportService) BranchSummary(filters dto.ReportFiltersDTO) ([]dto.BranchSummaryItemDTO, error) {
	return s.Repo.BranchSummary(filters)
}

func (s *ReportService) RoleDistribution() ([]dto.RoleDistributionItemDTO, error) {
	return s.Repo.RoleDistribution()
}

func (s *ReportService) UserGrowth(filters dto.ReportFiltersDTO) ([]dto.UserGrowthPointDTO, error) {
	if filters.GroupBy == "" {
		filters.GroupBy = "day"
	}

	return s.Repo.UserGrowth(filters)
}

func (s *ReportService) GlobalRoleDistribution(filters dto.ReportFiltersDTO) ([]dto.GlobalRoleItemDTO, error) {
	return s.Repo.GlobalRoleDistribution(filters)
}

func (s *ReportService) ActiveStatusDistribution(filters dto.ReportFiltersDTO) ([]dto.ActiveStatusItemDTO, error) {
	return s.Repo.ActiveStatusDistribution(filters)
}

func (s *ReportService) StaffSummary(filters dto.ReportFiltersDTO) (*dto.StaffSummaryResponse, error) {
	userTrends, err := s.UserGrowth(filters)
	if err != nil {
		return nil, err
	}

	globalRoles, err := s.GlobalRoleDistribution(filters)
	if err != nil {
		return nil, err
	}

	branchRoles, err := s.RoleDistribution()
	if err != nil {
		return nil, err
	}

	activeStatuses, err := s.ActiveStatusDistribution(filters)
	if err != nil {
		return nil, err
	}

	return &dto.StaffSummaryResponse{
		Message:        "staff summary fetched successfully",
		UserTrends:     userTrends,
		GlobalRoles:    globalRoles,
		BranchRoles:    branchRoles,
		ActiveStatuses: activeStatuses,
	}, nil
}

func (s *ReportService) BatchSummary(filters dto.ReportFiltersDTO) (*dto.BatchSummaryResponse, error) {
	totalBatches, err := s.TotalBatches(filters)
	if err != nil {
		return nil, err
	}

	statusCounts, err := s.Repo.BatchStatusSummary(filters)
	if err != nil {
		return nil, err
	}

	return &dto.BatchSummaryResponse{
		Message:      "batch summary fetched successfully",
		TotalBatches: totalBatches,
		StatusCounts: statusCounts,
	}, nil
}

func (s *ReportService) BranchTrends(filters dto.ReportFiltersDTO) ([]dto.BranchTrendSeriesDTO, error) {
	return s.Repo.BranchTrends(filters)
}

func (s *ReportService) GetBranches() ([]dto.BranchItemDTO, error) {
	return s.Repo.GetBranches()
}

func (s *ReportService) GetIngredientCategoryDaily(
	ctx context.Context,
	date time.Time,
) ([]dto.IngredientCategoryDailyDTO, error) {

	return s.Repo.GetIngredientCategoryDaily(ctx, date)
}

func (s *ReportService) ExportIngredientCategoryDailyExcel(
	ctx context.Context,
	date time.Time,
) ([]byte, error) {

	rows, err := s.Repo.GetIngredientCategoryDaily(ctx, date)
	if err != nil {
		return nil, err
	}

	file := excelize.NewFile()
	sheet := "Ingredient Category Daily"
	file.SetSheetName("Sheet1", sheet)

	headers := []string{
		"Used Date",
		"Category",
		"Base Unit",
		"Total Final",
		"Display Unit",
		"Display Final",
	}

	for index, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(index+1, 1)
		file.SetCellValue(sheet, cell, header)
	}

	for index, row := range rows {
		excelRow := index + 2

		file.SetCellValue(sheet, "A"+strconv.Itoa(excelRow), row.UsedDate)
		file.SetCellValue(sheet, "B"+strconv.Itoa(excelRow), row.CategoryName)
		file.SetCellValue(sheet, "C"+strconv.Itoa(excelRow), row.Unit)
		file.SetCellValue(sheet, "D"+strconv.Itoa(excelRow), exportFloat3(row.TotalFinalValue))
		file.SetCellValue(sheet, "E"+strconv.Itoa(excelRow), exportDisplayUnit(row.Unit))
		file.SetCellValue(sheet, "F"+strconv.Itoa(excelRow), exportFloat3(exportDisplayValue(row.Unit, row.TotalFinalValue)))
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *ReportService) ExportBatchDetailExcel(
	ctx context.Context,
	batchID int64,
) ([]byte, error) {

	batch, err := s.Repo.GetBatchDetailExport(ctx, batchID)
	if err != nil {
		return nil, err
	}

	file := excelize.NewFile()
	sheet := "Batch Detail"
	file.SetSheetName("Sheet1", sheet)

	// Header section
	file.SetCellValue(sheet, "A1", "Batch ID")
	file.SetCellValue(sheet, "B1", batch.BatchID)

	file.SetCellValue(sheet, "A2", "Recipe")
	file.SetCellValue(sheet, "B2", batch.RecipeName)

	file.SetCellValue(sheet, "A3", "Branch")
	file.SetCellValue(sheet, "B3", batch.BranchName)

	file.SetCellValue(sheet, "A4", "Created By")
	file.SetCellValue(sheet, "B4", batch.CreatedBy)

	file.SetCellValue(sheet, "A5", "Used Date")
	file.SetCellValue(sheet, "B5", batch.UsedDate)

	file.SetCellValue(sheet, "A6", "People")
	file.SetCellValue(sheet, "B6", batch.NPeople)

	file.SetCellValue(sheet, "A7", "Status")
	file.SetCellValue(sheet, "B7", batch.Status)

	// Table headers
	headers := []string{
		"Ingredient",
		"Base Unit",
		"Final",
		"Display Unit",
		"Display Final",
	}

	for index, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(index+1, 10)
		file.SetCellValue(sheet, cell, header)
	}

	// Table rows
	for index, item := range batch.Items {
		row := index + 11

		file.SetCellValue(sheet, "A"+strconv.Itoa(row), item.Ingredient)
		file.SetCellValue(sheet, "B"+strconv.Itoa(row), item.Unit)
		file.SetCellValue(sheet, "C"+strconv.Itoa(row), exportFloat3(item.FinalValue))
		file.SetCellValue(sheet, "D"+strconv.Itoa(row), exportDisplayUnit(item.Unit))
		file.SetCellValue(sheet, "E"+strconv.Itoa(row), exportFloat3(exportDisplayValue(item.Unit, item.FinalValue)))
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *ReportService) ExportBatchDetailPDF(
	ctx context.Context,
	batchID int64,
) ([]byte, error) {

	batch, err := s.Repo.GetBatchDetailExport(ctx, batchID)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(12, 12, 12)
	pdf.AddPage()

	logoPath := "assets/newco-logo.png"

	// Logo
	// Logo (left)
	pdf.ImageOptions(
		logoPath,
		12,
		10,
		45,
		0,
		false,
		gofpdf.ImageOptions{
			ImageType: "PNG",
			ReadDpi:   true,
		},
		0,
		"",
	)

	// Title (centered, lower than logo)
	pdf.SetY(22)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Consumption Detail Report", "", 1, "C", false, 0, "")

	// Divider line
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(12, 32, 198, 32)

	pdf.Ln(6)

	// Header info
	pdf.SetFont("Arial", "", 11)

	row := func(label1, value1, label2, value2 string) {
		x := pdf.GetX()
		y := pdf.GetY()

		pdf.SetXY(x, y)
		pdf.CellFormat(35, 6, label1, "", 0, "", false, 0, "")

		pdf.SetXY(x+35, y)
		pdf.MultiCell(55, 6, " "+value1, "", "", false)
		leftY := pdf.GetY()

		pdf.SetXY(x+95, y)
		pdf.CellFormat(35, 6, label2, "", 0, "", false, 0, "")

		pdf.SetXY(x+130, y)
		pdf.MultiCell(56, 6, " "+value2, "", "", false)
		rightY := pdf.GetY()

		nextY := leftY
		if rightY > nextY {
			nextY = rightY
		}

		pdf.SetY(nextY)
	}

	row("Batch ID:", strconv.FormatInt(batch.BatchID, 10), "Recipe:", batch.RecipeName)
	row("Branch:", batch.BranchName, "Created By:", batch.CreatedBy)
	row("Used Date:", batch.UsedDate, "People:", strconv.Itoa(batch.NPeople))
	row("Status:", batch.Status, "", "")
	pdf.Ln(8)

	const (
		colIngredient   = 78.0
		colBaseUnit     = 24.0
		colFinal        = 28.0
		colDisplayUnit  = 26.0
		colDisplayFinal = 30.0
	)

	// Table header helper
	writeTableHeader := func() {
		pdf.SetFont("Arial", "B", 11)

		pdf.CellFormat(colIngredient, 8, "Ingredient", "1", 0, "", false, 0, "")
		pdf.CellFormat(colBaseUnit, 8, "Base Unit", "1", 0, "C", false, 0, "")
		pdf.CellFormat(colFinal, 8, "Final", "1", 0, "R", false, 0, "")
		pdf.CellFormat(colDisplayUnit, 8, "Display Unit", "1", 0, "C", false, 0, "")
		pdf.CellFormat(colDisplayFinal, 8, "Display Final", "1", 1, "R", false, 0, "")

		pdf.SetFont("Arial", "", 10)
	}

	writeTableHeader()

	// Table rows
	pdf.SetFont("Arial", "", 10)

	for _, item := range batch.Items {
		if pdf.GetY() > 275 {
			pdf.AddPage()
			pdf.SetY(15)
			writeTableHeader()
		}

		pdf.CellFormat(colIngredient, 7, item.Ingredient, "1", 0, "", false, 0, "")
		pdf.CellFormat(colBaseUnit, 7, item.Unit, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colFinal, 7, exportFloat3(item.FinalValue), "1", 0, "R", false, 0, "")
		pdf.CellFormat(colDisplayUnit, 7, exportDisplayUnit(item.Unit), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colDisplayFinal, 7, exportFloat3(exportDisplayValue(item.Unit, item.FinalValue)), "1", 1, "R", false, 0, "")
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *ReportService) GetTopRecipeVariance(
	ctx context.Context,
	startDate string,
	endDate string,
	branchID string,
) (*dto.TopRecipeVarianceResponse, error) {

	items, err := s.Repo.GetTopRecipeVariance(
		ctx,
		startDate,
		endDate,
		branchID,
	)
	if err != nil {
		return nil, err
	}

	return &dto.TopRecipeVarianceResponse{
		Items: items,
	}, nil
}
