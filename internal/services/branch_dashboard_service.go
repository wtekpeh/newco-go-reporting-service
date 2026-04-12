package services

import (
	"context"
	"errors"

	"newco-go-reporting-service/internal/dto"
	"newco-go-reporting-service/internal/repositories"
)

var ErrNoBranchScope = errors.New("no branch scope available")

type BranchDashboardService interface {
	GetSummary(ctx context.Context, access *dto.AccessContext) (*dto.BranchDashboardSummaryResponse, error)
	GetBatchTrends(ctx context.Context, access *dto.AccessContext) ([]dto.BatchTrendPointDTO, error)
	GetRoleDistribution(ctx context.Context, access *dto.AccessContext) ([]dto.RoleDistributionItemDTO, error)
	GetRecentBatches(ctx context.Context, access *dto.AccessContext) ([]dto.RecentBatchItemDTO, error)
}

type branchDashboardService struct {
	repo repositories.BranchDashboardRepository
}

func NewBranchDashboardService(repo repositories.BranchDashboardRepository) BranchDashboardService {
	return &branchDashboardService{
		repo: repo,
	}
}

func (s *branchDashboardService) GetSummary(
	ctx context.Context,
	access *dto.AccessContext,
) (*dto.BranchDashboardSummaryResponse, error) {
	if access == nil || len(access.BranchIDs) == 0 {
		return nil, ErrNoBranchScope
	}

	branchID := access.BranchIDs[0]

	return s.repo.GetBranchSummary(ctx, branchID)
}

func (s *branchDashboardService) GetBatchTrends(
	ctx context.Context,
	access *dto.AccessContext,
) ([]dto.BatchTrendPointDTO, error) {

	if access == nil || len(access.BranchIDs) == 0 {
		return nil, ErrNoBranchScope
	}

	branchID := access.BranchIDs[0]

	return s.repo.GetBatchTrends(ctx, branchID)
}

func (s *branchDashboardService) GetRoleDistribution(
	ctx context.Context,
	access *dto.AccessContext,
) ([]dto.RoleDistributionItemDTO, error) {

	if access == nil || len(access.BranchIDs) == 0 {
		return nil, ErrNoBranchScope
	}

	branchID := access.BranchIDs[0]

	return s.repo.GetRoleDistribution(ctx, branchID)
}

func (s *branchDashboardService) GetRecentBatches(
	ctx context.Context,
	access *dto.AccessContext,
) ([]dto.RecentBatchItemDTO, error) {
	if access == nil || len(access.BranchIDs) == 0 {
		return nil, ErrNoBranchScope
	}

	branchID := access.BranchIDs[0]

	return s.repo.GetRecentBatches(ctx, branchID)
}
