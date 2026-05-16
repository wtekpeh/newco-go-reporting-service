package repositories

import (
	"context"
	"time"

	"newco-go-reporting-service/internal/dto"
	"newco-go-reporting-service/internal/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BranchDashboardRepository interface {
	GetBranchSummary(ctx context.Context, branchID int64) (*dto.BranchDashboardSummaryResponse, error)
	GetBatchTrends(ctx context.Context, branchID int64) ([]dto.BatchTrendPointDTO, error)
	GetDailyPlanTrends(ctx context.Context, branchID int64) ([]dto.BatchTrendPointDTO, error)
	GetRoleDistribution(ctx context.Context, branchID int64) ([]dto.RoleDistributionItemDTO, error)
	GetRecentBatches(ctx context.Context, branchID int64) ([]dto.RecentBatchItemDTO, error)
}

type branchDashboardRepository struct {
	pool *pgxpool.Pool
}

func NewBranchDashboardRepository(pool *pgxpool.Pool) BranchDashboardRepository {
	return &branchDashboardRepository{pool: pool}
}

func (r *branchDashboardRepository) GetBranchSummary(ctx context.Context, branchID int64) (*dto.BranchDashboardSummaryResponse, error) {
	var response dto.BranchDashboardSummaryResponse

	err := r.pool.QueryRow(ctx, queries.BranchDashboardSummaryQuery, branchID).Scan(
		&response.KPIs.BranchID,
		&response.KPIs.TotalStaff,
		&response.KPIs.TotalBatches,
		&response.KPIs.BatchesThisWeek,
		&response.KPIs.BatchesThisMonth,
		&response.KPIs.TotalDailyPlans,
		&response.KPIs.FinalizedDailyPlans,
		&response.KPIs.DraftDailyPlans,
	)
	if err != nil {
		return nil, err
	}

	recentDailyPlans, err := r.GetRecentDailyPlans(ctx, branchID)
	if err != nil {
		return nil, err
	}

	response.RecentDailyPlans = recentDailyPlans

	return &response, nil
}

func (r *branchDashboardRepository) GetRecentDailyPlans(
	ctx context.Context,
	branchID int64,
) ([]dto.BranchRecentDailyPlanItem, error) {
	rows, err := r.pool.Query(
		ctx,
		queries.BranchRecentDailyPlansQuery,
		branchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []dto.BranchRecentDailyPlanItem{}

	for rows.Next() {
		var item dto.BranchRecentDailyPlanItem
		var usedDate time.Time
		var createdAt time.Time

		err := rows.Scan(
			&item.DailyPlanID,
			&usedDate,
			&item.Status,
			&item.CreatedBy,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		item.UsedDate = usedDate.Format("2006-01-02")
		item.CreatedAt = createdAt.Format(time.RFC3339)

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *branchDashboardRepository) GetBatchTrends(
	ctx context.Context,
	branchID int64,
) ([]dto.BatchTrendPointDTO, error) {

	rows, err := r.pool.Query(
		ctx,
		queries.BranchBatchTrendsQuery,
		branchID,
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

func (r *branchDashboardRepository) GetRoleDistribution(
	ctx context.Context,
	branchID int64,
) ([]dto.RoleDistributionItemDTO, error) {

	rows, err := r.pool.Query(
		ctx,
		queries.BranchRoleDistributionQuery,
		branchID,
	)
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

func (r *branchDashboardRepository) GetRecentBatches(
	ctx context.Context,
	branchID int64,
) ([]dto.RecentBatchItemDTO, error) {
	rows, err := r.pool.Query(ctx, queries.BranchRecentBatchesQuery, branchID)
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

func (r *branchDashboardRepository) GetDailyPlanTrends(
	ctx context.Context,
	branchID int64,
) ([]dto.BatchTrendPointDTO, error) {
	rows, err := r.pool.Query(
		ctx,
		queries.BranchDailyPlanTrendsQuery,
		branchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []dto.BatchTrendPointDTO{}

	for rows.Next() {
		var item dto.BatchTrendPointDTO

		err := rows.Scan(
			&item.Label,
			&item.Count,
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
