package services

import "newco-go-reporting-service/internal/repositories"

type ReportService struct {
	Repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{
		Repo: repo,
	}
}

func (s *ReportService) TotalUsers() (int, error) {
	return s.Repo.TotalUsers()
}

func (s *ReportService) TotalActiveUsers() (int, error) {
	return s.Repo.TotalActiveUsers()
}
