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
