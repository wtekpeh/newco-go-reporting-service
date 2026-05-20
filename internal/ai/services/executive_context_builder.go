package services

import (
	"context"
	"encoding/json"

	aidto "newco-go-reporting-service/internal/ai/dto"
	reportdto "newco-go-reporting-service/internal/dto"
	reportservices "newco-go-reporting-service/internal/services"
)

type ExecutiveContextBuilder struct {
	ReportService *reportservices.ReportService
}

func NewExecutiveContextBuilder(
	reportService *reportservices.ReportService,
) *ExecutiveContextBuilder {
	return &ExecutiveContextBuilder{
		ReportService: reportService,
	}
}

type ExecutiveAIContext struct {
	ReportingPeriod ExecutiveAIReportingPeriod `json:"reporting_period"`
	KPIs            ExecutiveAIKpis            `json:"kpis"`
	Highlights      ExecutiveAIHighlights      `json:"highlights"`
}

type ExecutiveAIReportingPeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type ExecutiveAIKpis struct {
	TotalUsers       int `json:"total_users"`
	TotalActiveUsers int `json:"total_active_users"`
	TotalBranches    int `json:"total_branches"`

	TotalBatches int `json:"total_batches"`
}

type ExecutiveAIHighlights struct {
	MostActiveBranch        any `json:"most_active_branch"`
	LargestBranch           any `json:"largest_branch"`
	MostUsedRecipe          any `json:"most_used_recipe"`
	AverageBatchesPerBranch any `json:"average_batches_per_branch"`
	PeakBatchDay            any `json:"peak_batch_day"`
}

func (b *ExecutiveContextBuilder) BuildExecutiveSummaryContext(
	ctx context.Context,
	request aidto.ExecutiveSummaryRequest,
) (string, error) {

	filters := reportdto.ReportFiltersDTO{
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
	}

	if request.BranchID != nil {
		filters.BranchID = string(rune(*request.BranchID))
	}

	totalUsers, err := b.ReportService.TotalUsers(filters)
	if err != nil {
		return "", err
	}

	totalActiveUsers, err := b.ReportService.TotalActiveUsers(filters)
	if err != nil {
		return "", err
	}

	totalBranches, err := b.ReportService.TotalBranches()
	if err != nil {
		return "", err
	}

	totalBatches, err := b.ReportService.TotalBatches(filters)
	if err != nil {
		return "", err
	}

	mostActiveBranch, err := b.ReportService.MostActiveBranch(filters)
	if err != nil {
		return "", err
	}

	largestBranch, err := b.ReportService.LargestBranch(filters)
	if err != nil {
		return "", err
	}

	mostUsedRecipe, err := b.ReportService.MostUsedRecipe(filters)
	if err != nil {
		return "", err
	}

	averageBatchesPerBranch, err := b.ReportService.AverageBatchesPerBranch(filters)
	if err != nil {
		return "", err
	}

	peakBatchDay, err := b.ReportService.PeakBatchDay(filters)
	if err != nil {
		return "", err
	}

	aiContext := ExecutiveAIContext{
		ReportingPeriod: ExecutiveAIReportingPeriod{
			StartDate: request.StartDate,
			EndDate:   request.EndDate,
		},
		KPIs: ExecutiveAIKpis{
			TotalUsers:       totalUsers,
			TotalActiveUsers: totalActiveUsers,
			TotalBranches:    totalBranches,
			TotalBatches:     totalBatches,
		},
		Highlights: ExecutiveAIHighlights{
			MostActiveBranch:        mostActiveBranch,
			LargestBranch:           largestBranch,
			MostUsedRecipe:          mostUsedRecipe,
			AverageBatchesPerBranch: averageBatchesPerBranch,
			PeakBatchDay:            peakBatchDay,
		},
	}

	jsonBytes, err := json.MarshalIndent(aiContext, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
