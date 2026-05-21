package services

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

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
	ReportingPeriod      ExecutiveAIReportingPeriod      `json:"reporting_period"`
	KPIs                 ExecutiveAIKpis                 `json:"kpis"`
	Highlights           ExecutiveAIHighlights           `json:"highlights"`
	IngredientCategories []ExecutiveAIIngredientCategory `json:"ingredient_categories"`
	TopRecipeVariance    []ExecutiveAITopRecipeVariance  `json:"top_recipe_variance"`
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

	TotalDailyPlans     int `json:"total_daily_plans"`
	FinalizedDailyPlans int `json:"finalized_daily_plans"`
	DraftDailyPlans     int `json:"draft_daily_plans"`
}

type ExecutiveAIHighlights struct {
	MostActiveBranch        any `json:"most_active_branch"`
	LargestBranch           any `json:"largest_branch"`
	MostUsedRecipe          any `json:"most_used_recipe"`
	AverageBatchesPerBranch any `json:"average_batches_per_branch"`
	PeakBatchDay            any `json:"peak_batch_day"`
}

type ExecutiveAIIngredientCategory struct {
	CategoryName string  `json:"category_name"`
	Unit         string  `json:"unit"`
	TotalFinal   float64 `json:"total_final"`
	TotalActual  float64 `json:"total_actual"`
}

type ExecutiveAITopRecipeVariance struct {
	RecipeName      string  `json:"recipe_name"`
	VariancePercent float64 `json:"variance_percent"`
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

	totalDailyPlans, err := b.ReportService.TotalDailyPlans(filters)
	if err != nil {
		return "", err
	}

	finalizedDailyPlans, err := b.ReportService.FinalizedDailyPlans(filters)
	if err != nil {
		return "", err
	}

	draftDailyPlans, err := b.ReportService.DraftDailyPlans(filters)
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

	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return "", err
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return "", err
	}

	ingredientCategoryReport, err := b.ReportService.GetIngredientCategoryDaily(
		ctx,
		startDate,
		endDate,
	)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}

	ingredientCategories := []ExecutiveAIIngredientCategory{}

	topRecipeVariance := []ExecutiveAITopRecipeVariance{}

	for _, item := range ingredientCategoryReport {
		ingredientCategories = append(
			ingredientCategories,
			ExecutiveAIIngredientCategory{
				CategoryName: item.CategoryName,
				Unit:         item.Unit,
				TotalFinal:   item.TotalFinalValue,
				TotalActual:  item.TotalActualValue,
			},
		)
	}

	branchID := ""

	if request.BranchID != nil {
		branchID = strconv.FormatInt(*request.BranchID, 10)
	}

	topRecipeVarianceResult, err := b.ReportService.GetTopRecipeVariance(
		ctx,
		request.StartDate,
		request.EndDate,
		branchID,
	)
	if err != nil {
		return "", err
	}

	for _, item := range topRecipeVarianceResult.Items {
		topRecipeVariance = append(
			topRecipeVariance,
			ExecutiveAITopRecipeVariance{
				RecipeName:      item.RecipeName,
				VariancePercent: item.AverageVarianceG,
			},
		)
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

			TotalDailyPlans:     totalDailyPlans,
			FinalizedDailyPlans: finalizedDailyPlans,
			DraftDailyPlans:     draftDailyPlans,
		},
		Highlights: ExecutiveAIHighlights{
			MostActiveBranch:        mostActiveBranch,
			LargestBranch:           largestBranch,
			MostUsedRecipe:          mostUsedRecipe,
			AverageBatchesPerBranch: averageBatchesPerBranch,
			PeakBatchDay:            peakBatchDay,
		},
		IngredientCategories: ingredientCategories,
		TopRecipeVariance:    topRecipeVariance,
	}

	jsonBytes, err := json.MarshalIndent(aiContext, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
