package services

import (
	"context"

	"newco-go-reporting-service/internal/dto"
	"newco-go-reporting-service/internal/repositories"
)

type AccessService interface {
	ResolveAccessContext(ctx context.Context, keycloakSub string) (*dto.AccessContext, error)
	IsExecutive(access *dto.AccessContext) bool
	IsBranchManager(access *dto.AccessContext) bool
}

type accessService struct {
	accessRepo       repositories.AccessRepository
	branchAccessRepo repositories.BranchAccessRepository
}

func NewAccessService(
	accessRepo repositories.AccessRepository,
	branchAccessRepo repositories.BranchAccessRepository,
) AccessService {
	return &accessService{
		accessRepo:       accessRepo,
		branchAccessRepo: branchAccessRepo,
	}
}

func (s *accessService) ResolveAccessContext(ctx context.Context, keycloakSub string) (*dto.AccessContext, error) {
	record, err := s.accessRepo.GetActiveStaffByKeycloakSub(ctx, keycloakSub)
	if err != nil {
		return nil, err
	}

	branchRoles, err := s.branchAccessRepo.GetActiveBranchRolesByStaffProfileID(ctx, record.StaffProfileID)
	if err != nil {
		return nil, err
	}

	branchIDs := make([]int64, 0)
	seen := make(map[int64]struct{})
	isBranchManager := false

	for _, role := range branchRoles {
		if role.Role == "branch_manager" {
			isBranchManager = true
			if _, exists := seen[role.BranchID]; !exists {
				seen[role.BranchID] = struct{}{}
				branchIDs = append(branchIDs, role.BranchID)
			}
		}
	}

	access := &dto.AccessContext{
		StaffProfileID:  record.StaffProfileID,
		KeycloakSub:     record.KeycloakSub,
		Email:           record.Email,
		Username:        record.Username,
		FullName:        record.FullName,
		GlobalRole:      record.GlobalRole,
		IsExecutive:     s.isExecutiveRole(record.GlobalRole),
		BranchIDs:       branchIDs,
		IsBranchManager: isBranchManager,
	}

	return access, nil
}

func (s *accessService) IsExecutive(access *dto.AccessContext) bool {
	if access == nil {
		return false
	}
	return s.isExecutiveRole(access.GlobalRole)
}

func (s *accessService) IsBranchManager(access *dto.AccessContext) bool {
	if access == nil {
		return false
	}
	return access.IsBranchManager && len(access.BranchIDs) > 0
}

func (s *accessService) isExecutiveRole(globalRole string) bool {
	return globalRole == "boss" || globalRole == "managing_director"
}
