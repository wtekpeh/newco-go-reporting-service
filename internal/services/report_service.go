package services

import (
	"context"
	"newco-go-reporting-service/internal/dto"
	"newco-go-reporting-service/internal/repositories"
	"time"
)

type ReportService struct {
	Repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{
		Repo: repo,
	}
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
