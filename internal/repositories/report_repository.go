package repositories

import (
	"context"
	"errors"
	"newco-go-reporting-service/internal/dto"
	"newco-go-reporting-service/internal/queries"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	DB *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{
		DB: db,
	}
}

func (r *ReportRepository) TotalUsers(filters dto.ReportFiltersDTO) (int, error) {
	var count int

	query := queries.TotalUsersQuery

	err := r.DB.QueryRow(context.Background(), query, filters.BranchID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepository) TotalActiveUsers(filters dto.ReportFiltersDTO) (int, error) {
	var count int

	query := queries.TotalActiveUsersQuery

	err := r.DB.QueryRow(context.Background(), query, filters.BranchID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepository) TotalBranches() (int, error) {
	var count int

	query := queries.TotalBranchesQuery

	err := r.DB.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepository) TotalBatches(filters dto.ReportFiltersDTO) (int, error) {
	var count int

	query := queries.TotalBatchesQuery

	err := r.DB.QueryRow(context.Background(), query, filters.BranchID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepository) BatchesThisWeek(filters dto.ReportFiltersDTO) (int, error) {
	var count int

	query := queries.BatchesThisWeekQuery

	err := r.DB.QueryRow(context.Background(), query, filters.BranchID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepository) BatchesThisMonth(filters dto.ReportFiltersDTO) (int, error) {
	var count int

	query := queries.BatchesThisMonthQuery

	err := r.DB.QueryRow(context.Background(), query, filters.BranchID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepository) MostActiveBranch(filters dto.ReportFiltersDTO) (*dto.BranchActivityDTO, error) {
	var result dto.BranchActivityDTO

	query := queries.MostActiveBranchQuery

	err := r.DB.QueryRow(context.Background(), query, filters.BranchID).Scan(&result.BranchName, &result.BatchCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *ReportRepository) LargestBranch(filters dto.ReportFiltersDTO) (*dto.BranchStaffDTO, error) {
	var result dto.BranchStaffDTO

	query := queries.LargestBranchQuery

	err := r.DB.QueryRow(context.Background(), query, filters.BranchID).Scan(&result.BranchName, &result.StaffCount)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *ReportRepository) MostUsedRecipe(filters dto.ReportFiltersDTO) (*dto.RecipeUsageDTO, error) {
	var result dto.RecipeUsageDTO

	query := queries.MostUsedRecipeQuery

	err := r.DB.QueryRow(context.Background(), query, filters.BranchID).Scan(&result.RecipeName, &result.BatchCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *ReportRepository) AverageBatchesPerBranch(filters dto.ReportFiltersDTO) (*dto.AverageMetricDTO, error) {
	var result dto.AverageMetricDTO

	query := queries.AverageBatchesPerBranchQuery

	err := r.DB.QueryRow(context.Background(), query, filters.BranchID).Scan(&result.Value)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *ReportRepository) PeakBatchDay(filters dto.ReportFiltersDTO) (*dto.PeakBatchDayDTO, error) {
	var result dto.PeakBatchDayDTO

	query := queries.PeakBatchDayQuery

	err := r.DB.QueryRow(context.Background(), query, filters.BranchID).Scan(&result.DayName, &result.BatchCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *ReportRepository) RecentBatches(filters dto.ReportFiltersDTO) ([]dto.RecentBatchItemDTO, error) {
	rows, err := r.DB.Query(
		context.Background(),
		queries.RecentBatchesQuery,
		filters.StartDate,
		filters.EndDate,
		filters.BranchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.RecentBatchItemDTO

	for rows.Next() {
		var item dto.RecentBatchItemDTO
		var createdAt time.Time

		err := rows.Scan(
			&item.BatchID,
			&item.RecipeName,
			&item.BranchName,
			&item.CreatedBy,
			&item.NPeople,
			&item.Status,
			&item.ProteinType,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ReportRepository) BatchTrends(filters dto.ReportFiltersDTO) ([]dto.BatchTrendPointDTO, error) {
	rows, err := r.DB.Query(
		context.Background(),
		queries.BatchTrendsQuery,
		filters.StartDate,
		filters.EndDate,
		filters.BranchID,
		filters.GroupBy,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.BatchTrendPointDTO

	for rows.Next() {
		var item dto.BatchTrendPointDTO

		err := rows.Scan(&item.Label, &item.Count)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ReportRepository) BranchSummary(filters dto.ReportFiltersDTO) ([]dto.BranchSummaryItemDTO, error) {
	rows, err := r.DB.Query(context.Background(), queries.BranchSummaryQuery, filters.BranchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.BranchSummaryItemDTO

	for rows.Next() {
		var item dto.BranchSummaryItemDTO

		err := rows.Scan(
			&item.BranchID,
			&item.BranchName,
			&item.StaffCount,
			&item.TotalBatches,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ReportRepository) RoleDistribution() ([]dto.RoleDistributionItemDTO, error) {
	rows, err := r.DB.Query(context.Background(), queries.RoleDistributionQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.RoleDistributionItemDTO

	for rows.Next() {
		var item dto.RoleDistributionItemDTO

		err := rows.Scan(&item.Role, &item.Count)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ReportRepository) UserGrowth(filters dto.ReportFiltersDTO) ([]dto.UserGrowthPointDTO, error) {
	rows, err := r.DB.Query(
		context.Background(),
		queries.UserGrowthQuery,
		filters.StartDate,
		filters.EndDate,
		filters.BranchID,
		filters.GroupBy,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.UserGrowthPointDTO

	for rows.Next() {
		var item dto.UserGrowthPointDTO

		err := rows.Scan(&item.Label, &item.Count)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ReportRepository) GlobalRoleDistribution(filters dto.ReportFiltersDTO) ([]dto.GlobalRoleItemDTO, error) {
	rows, err := r.DB.Query(context.Background(), queries.GlobalRoleDistributionQuery, filters.BranchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.GlobalRoleItemDTO

	for rows.Next() {
		var item dto.GlobalRoleItemDTO

		err := rows.Scan(&item.Role, &item.Count)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ReportRepository) ActiveStatusDistribution(filters dto.ReportFiltersDTO) ([]dto.ActiveStatusItemDTO, error) {
	rows, err := r.DB.Query(context.Background(), queries.ActiveStatusDistributionQuery, filters.BranchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.ActiveStatusItemDTO

	for rows.Next() {
		var item dto.ActiveStatusItemDTO

		err := rows.Scan(&item.Status, &item.Count)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ReportRepository) BatchStatusSummary(filters dto.ReportFiltersDTO) ([]dto.BatchStatusItemDTO, error) {
	rows, err := r.DB.Query(context.Background(), queries.BatchStatusSummaryQuery, filters.BranchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.BatchStatusItemDTO

	for rows.Next() {
		var item dto.BatchStatusItemDTO

		err := rows.Scan(&item.Status, &item.Count)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ReportRepository) BranchTrends(filters dto.ReportFiltersDTO) ([]dto.BranchTrendSeriesDTO, error) {
	rows, err := r.DB.Query(context.Background(), queries.BranchTrendsFlatQuery, filters.BranchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seriesMap := make(map[int]*dto.BranchTrendSeriesDTO)
	var order []int

	for rows.Next() {
		var branchID int
		var branchName string
		var point dto.BranchTrendPointDTO

		err := rows.Scan(&branchID, &branchName, &point.Label, &point.Count)
		if err != nil {
			return nil, err
		}

		if _, exists := seriesMap[branchID]; !exists {
			seriesMap[branchID] = &dto.BranchTrendSeriesDTO{
				BranchID:   branchID,
				BranchName: branchName,
				Points:     []dto.BranchTrendPointDTO{},
			}
			order = append(order, branchID)
		}

		seriesMap[branchID].Points = append(seriesMap[branchID].Points, point)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var series []dto.BranchTrendSeriesDTO
	for _, id := range order {
		series = append(series, *seriesMap[id])
	}

	return series, nil
}

func (r *ReportRepository) GetBranches() ([]dto.BranchItemDTO, error) {
	rows, err := r.DB.Query(context.Background(), `
		SELECT id, name
		FROM accounts_branch
		WHERE is_active = TRUE
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.BranchItemDTO

	for rows.Next() {
		var item dto.BranchItemDTO

		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *ReportRepository) GetBranchByName(
	name string,
) (dto.BranchLookupDTO, error) {

	query := `
		SELECT
			id,
			name
		FROM accounts_branch
		WHERE is_active = TRUE
		  AND name ILIKE '%' || $1 || '%'
		ORDER BY name ASC
		LIMIT 1
	`

	var branch dto.BranchLookupDTO

	err := r.DB.QueryRow(
		context.Background(),
		query,
		name,
	).Scan(
		&branch.ID,
		&branch.Name,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.BranchLookupDTO{}, nil
		}

		return dto.BranchLookupDTO{}, err
	}

	return branch, nil
}

func (r *ReportRepository) GetIngredientCategoryDaily(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]dto.IngredientCategoryDailyDTO, error) {

	rows, err := r.DB.Query(
		ctx,
		queries.IngredientCategoryDailyQuery,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []dto.IngredientCategoryDailyDTO{}

	for rows.Next() {
		var item dto.IngredientCategoryDailyDTO

		err := rows.Scan(
			&item.UsedDate,
			&item.CategoryID,
			&item.CategoryName,
			&item.Unit,
			&item.TotalFinalValue,
			&item.TotalActualValue,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return results, nil
}

func (r *ReportRepository) GetBatchDetailExport(
	ctx context.Context,
	batchID int64,
) (*dto.BatchDetailExportDTO, error) {

	var result dto.BatchDetailExportDTO

	err := r.DB.QueryRow(ctx, queries.BatchDetailExportHeaderQuery, batchID).Scan(
		&result.BatchID,
		&result.RecipeName,
		&result.BranchName,
		&result.CreatedBy,
		&result.UsedDate,
		&result.NPeople,
		&result.Status,
	)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, queries.BatchDetailExportItemsQuery, batchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []dto.BatchDetailItemDTO{}

	for rows.Next() {
		var item dto.BatchDetailItemDTO

		err := rows.Scan(
			&item.Ingredient,
			&item.Unit,
			&item.FinalValue,
			&item.ActualValue,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	result.Items = items

	return &result, nil
}

func (r *ReportRepository) GetTopRecipeVariance(
	ctx context.Context,
	startDate string,
	endDate string,
	branchID string,
) ([]dto.TopRecipeVarianceItem, error) {
	rows, err := r.DB.Query(
		ctx,
		queries.TopRecipeVarianceQuery,
		startDate,
		endDate,
		branchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []dto.TopRecipeVarianceItem{}

	for rows.Next() {
		var item dto.TopRecipeVarianceItem

		err := rows.Scan(
			&item.RecipeID,
			&item.RecipeName,
			&item.AverageVarianceG,
			&item.TotalBatches,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *ReportRepository) TotalDailyPlans(filters dto.ReportFiltersDTO) (int, error) {
	var count int

	err := r.DB.QueryRow(
		context.Background(),
		queries.TotalDailyPlansQuery,
		filters.BranchID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepository) FinalizedDailyPlans(filters dto.ReportFiltersDTO) (int, error) {
	var count int

	err := r.DB.QueryRow(
		context.Background(),
		queries.FinalizedDailyPlansQuery,
		filters.BranchID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepository) DraftDailyPlans(filters dto.ReportFiltersDTO) (int, error) {
	var count int

	err := r.DB.QueryRow(
		context.Background(),
		queries.DraftDailyPlansQuery,
		filters.BranchID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReportRepository) DailyPlanTrends(filters dto.ReportFiltersDTO) ([]dto.DailyPlanTrendPointDTO, error) {
	rows, err := r.DB.Query(
		context.Background(),
		queries.DailyPlanTrendsQuery,
		filters.StartDate,
		filters.EndDate,
		filters.BranchID,
		filters.GroupBy,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.DailyPlanTrendPointDTO

	for rows.Next() {
		var item dto.DailyPlanTrendPointDTO

		err := rows.Scan(&item.Label, &item.Count)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ReportRepository) RecentDailyPlans(filters dto.ReportFiltersDTO) ([]dto.RecentDailyPlanItemDTO, error) {
	rows, err := r.DB.Query(
		context.Background(),
		queries.RecentDailyPlansQuery,
		filters.StartDate,
		filters.EndDate,
		filters.BranchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []dto.RecentDailyPlanItemDTO

	for rows.Next() {
		var item dto.RecentDailyPlanItemDTO
		var planDate time.Time
		var createdAt time.Time

		err := rows.Scan(
			&item.DailyPlanID,
			&item.BranchName,
			&item.CreatedBy,
			&planDate,
			&item.Status,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		item.PlanDate = planDate.Format("2006-01-02")
		item.CreatedAt = createdAt.Format(time.RFC3339)

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ReportRepository) DailyPlanRequisitionExport(
	ctx context.Context,
	dailyPlanID int64,
) (*dto.DailyPlanRequisitionExportDTO, error) {

	rows, err := r.DB.Query(
		ctx,
		queries.DailyPlanRequisitionPDFQuery,
		dailyPlanID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result dto.DailyPlanRequisitionExportDTO
	items := []dto.DailyPlanRequisitionItemDTO{}

	for rows.Next() {
		var item dto.DailyPlanRequisitionItemDTO
		var usedDate time.Time
		var adjustedTotalKG float64

		err := rows.Scan(
			&result.DailyPlanID,
			&usedDate,
			&result.Status,
			&result.BranchName,
			&result.CreatedBy,
			&item.Ingredient,
			&item.Group,
			&item.Quantity,
			&adjustedTotalKG,
			&item.Unit,
		)
		if err != nil {
			return nil, err
		}

		result.UsedDate = usedDate.Format("2006-01-02")

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result.Items = items

	return &result, nil
}

func (r *ReportRepository) AIBranchPerformance(
	filters dto.ReportFiltersDTO,
) ([]dto.BranchSummaryItemDTO, error) {

	rows, err := r.DB.Query(
		context.Background(),
		queries.AIBranchPerformanceQuery,
		filters.StartDate,
		filters.EndDate,
		filters.BranchID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := []dto.BranchSummaryItemDTO{}

	for rows.Next() {

		var item dto.BranchSummaryItemDTO

		err := rows.Scan(
			&item.BranchID,
			&item.BranchName,
			&item.TotalBatches,
			&item.StaffCount,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *ReportRepository) DailyPlanSummary(
	filters dto.ReportFiltersDTO,
) ([]map[string]interface{}, error) {

	rows, err := r.DB.Query(
		context.Background(),
		queries.DailyPlanSummaryQuery,
		filters.StartDate,
		filters.EndDate,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := []map[string]interface{}{}

	for rows.Next() {

		var status string
		var count int64

		err := rows.Scan(
			&status,
			&count,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"status": status,
			"count":  count,
		})
	}

	return results, nil
}

func (r *ReportRepository) GetAIPlanningRiskSummary(
	filters dto.ReportFiltersDTO,
) (dto.AIPlanningRiskSummaryDTO, error) {

	var summary dto.AIPlanningRiskSummaryDTO

	err := r.DB.QueryRow(
		context.Background(),
		queries.AIPlanningRiskSummaryQuery,
		filters.StartDate,
		filters.EndDate,
	).Scan(
		&summary.TotalPlans,
		&summary.DraftPlans,
		&summary.FinalizedPlans,
		&summary.PlansMissingActuals,
	)
	if err != nil {
		return summary, err
	}

	// Risk evaluation
	if summary.DraftPlans > 0 {
		summary.RiskLevel = "medium"

		summary.RiskReasons = append(
			summary.RiskReasons,
			"There are draft plans pending finalization.",
		)

		summary.ManagementAttention = append(
			summary.ManagementAttention,
			"Review and finalize pending operational plans.",
		)
	}

	if summary.PlansMissingActuals > 0 {
		summary.RiskLevel = "medium"

		summary.RiskReasons = append(
			summary.RiskReasons,
			"Some finalized plans are missing actual execution values.",
		)

		summary.ManagementAttention = append(
			summary.ManagementAttention,
			"Monitor execution reporting completeness.",
		)
	}

	if summary.DraftPlans == 0 &&
		summary.PlansMissingActuals == 0 {

		summary.RiskLevel = "low"

		summary.RiskReasons = append(
			summary.RiskReasons,
			"Operational planning appears fully finalized.",
		)

		summary.ManagementAttention = append(
			summary.ManagementAttention,
			"Continue monitoring execution readiness.",
		)
	}

	return summary, nil
}

func (r *ReportRepository) GetAIIngredientVarianceRisk(
	filters dto.ReportFiltersDTO,
) (dto.AIIngredientVarianceRiskDTO, error) {

	rows, err := r.DB.Query(
		context.Background(),
		queries.AIIngredientVarianceRiskQuery,
		filters.StartDate,
		filters.EndDate,
		filters.BranchID,
	)
	if err != nil {
		return dto.AIIngredientVarianceRiskDTO{}, err
	}
	defer rows.Close()

	result := dto.AIIngredientVarianceRiskDTO{
		Items:               []dto.AIIngredientVarianceRiskItemDTO{},
		HighestRiskLevel:    "low",
		RiskReasons:         []string{},
		ManagementAttention: []string{},
	}

	for rows.Next() {
		var item dto.AIIngredientVarianceRiskItemDTO

		err := rows.Scan(
			&item.Ingredient,
			&item.Unit,
			&item.TotalPlannedValue,
			&item.TotalActualValue,
			&item.VarianceValue,
			&item.VariancePercent,
			&item.TotalConsumptions,
		)
		if err != nil {
			return result, err
		}

		if item.VariancePercent >= 30 {
			item.RiskLevel = "high"
			item.RiskReason = "Actual usage is significantly different from planned usage."
			result.HighestRiskLevel = "high"
		} else if item.VariancePercent >= 15 {
			item.RiskLevel = "medium"
			item.RiskReason = "Actual usage is moderately different from planned usage."

			if result.HighestRiskLevel != "high" {
				result.HighestRiskLevel = "medium"
			}
		} else {
			item.RiskLevel = "low"
			item.RiskReason = "Actual usage is close to planned usage."
		}

		result.Items = append(result.Items, item)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	if len(result.RiskReasons) == 0 {
		for _, item := range result.Items {
			if item.RiskLevel == "high" || item.RiskLevel == "medium" {
				result.RiskReasons = append(
					result.RiskReasons,
					item.Ingredient+" has noticeable planned versus actual usage variance.",
				)
			}
		}
	}

	if len(result.ManagementAttention) == 0 {
		if result.HighestRiskLevel == "high" {
			result.ManagementAttention = append(
				result.ManagementAttention,
				"Review high-variance ingredients and confirm whether recipes, planning assumptions, or execution records need adjustment.",
			)
		} else if result.HighestRiskLevel == "medium" {
			result.ManagementAttention = append(
				result.ManagementAttention,
				"Monitor medium-variance ingredients and compare them against recent operational patterns.",
			)
		} else {
			result.ManagementAttention = append(
				result.ManagementAttention,
				"Ingredient usage appears close to planned values.",
			)
		}
	}

	return result, nil
}
