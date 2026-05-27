package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
) (string, ExecutiveAIContext, error) {

	filters := reportdto.ReportFiltersDTO{
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
	}

	if request.BranchID != nil {
		filters.BranchID = string(rune(*request.BranchID))
	}

	totalUsers, err := b.ReportService.TotalUsers(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	totalActiveUsers, err := b.ReportService.TotalActiveUsers(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	totalBranches, err := b.ReportService.TotalBranches()
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	totalBatches, err := b.ReportService.TotalBatches(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	totalDailyPlans, err := b.ReportService.TotalDailyPlans(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	finalizedDailyPlans, err := b.ReportService.FinalizedDailyPlans(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	draftDailyPlans, err := b.ReportService.DraftDailyPlans(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	mostActiveBranch, err := b.ReportService.MostActiveBranch(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	largestBranch, err := b.ReportService.LargestBranch(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	mostUsedRecipe, err := b.ReportService.MostUsedRecipe(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	averageBatchesPerBranch, err := b.ReportService.AverageBatchesPerBranch(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	peakBatchDay, err := b.ReportService.PeakBatchDay(filters)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}

	var startDate time.Time
	var endDate time.Time

	if request.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", request.StartDate)
		if err != nil {
			return "", ExecutiveAIContext{}, err
		}
	}

	if request.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", request.EndDate)
		if err != nil {
			return "", ExecutiveAIContext{}, err
		}
	}

	ingredientCategoryReport, err := b.ReportService.GetIngredientCategoryDaily(
		ctx,
		startDate,
		endDate,
	)
	if err != nil {
		return "", ExecutiveAIContext{}, err
	}
	if err != nil {
		return "", ExecutiveAIContext{}, err
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
		return "", ExecutiveAIContext{}, err
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
		return "", ExecutiveAIContext{}, err
	}

	jsonString := string(jsonBytes)

	jsonString = strings.ReplaceAll(
		jsonString,
		"branches",
		"sites",
	)

	jsonString = strings.ReplaceAll(
		jsonString,
		"Branches",
		"Sites",
	)

	jsonString = strings.ReplaceAll(
		jsonString,
		"branch",
		"site",
	)

	jsonString = strings.ReplaceAll(
		jsonString,
		"Branch",
		"Site",
	)

	jsonString = strings.ReplaceAll(
		jsonString,
		"batches",
		"consumptions",
	)

	jsonString = strings.ReplaceAll(
		jsonString,
		"Batches",
		"Consumptions",
	)

	jsonString = strings.ReplaceAll(
		jsonString,
		"batch",
		"consumption",
	)

	jsonString = strings.ReplaceAll(
		jsonString,
		"Batch",
		"Consumption",
	)

	return jsonString, aiContext, nil
}

func BuildExecutiveSummaryNarrative(
	context ExecutiveAIContext,
) string {

	summary := fmt.Sprintf(
		"Operations remain active across %d sites with %d total consumptions.",
		context.KPIs.TotalBranches,
		context.KPIs.TotalBatches,
	)

	summary += "\n\nKey highlights:\n"

	summary += fmt.Sprintf(
		"- Daily planning is complete with %d finalized plans and %d drafts currently pending.\n",
		context.KPIs.FinalizedDailyPlans,
		context.KPIs.DraftDailyPlans,
	)

	mostActiveSiteBytes, _ := json.Marshal(
		context.Highlights.MostActiveBranch,
	)

	var mostActiveSite map[string]any

	json.Unmarshal(
		mostActiveSiteBytes,
		&mostActiveSite,
	)

	siteName := mostActiveSite["branch_name"]
	consumptionCount := mostActiveSite["batch_count"]

	if siteName != nil {
		summary += fmt.Sprintf(
			"- %v remains the most active site with %v consumptions.\n",
			siteName,
			consumptionCount,
		)
	}

	largestSiteBytes, _ := json.Marshal(
		context.Highlights.LargestBranch,
	)

	var largestSite map[string]any

	json.Unmarshal(
		largestSiteBytes,
		&largestSite,
	)

	largestSiteName := largestSite["branch_name"]
	staffCount := largestSite["staff_count"]

	if largestSiteName != nil {
		summary += fmt.Sprintf(
			"- %v is currently the largest operational site with %v staff.\n",
			largestSiteName,
			staffCount,
		)
	}

	mostUsedRecipeBytes, _ := json.Marshal(
		context.Highlights.MostUsedRecipe,
	)

	var mostUsedRecipe map[string]any

	json.Unmarshal(
		mostUsedRecipeBytes,
		&mostUsedRecipe,
	)

	recipeName := mostUsedRecipe["recipe_name"]
	recipeConsumptionCount := mostUsedRecipe["batch_count"]

	if recipeName != nil {
		summary += fmt.Sprintf(
			"- %v is the most frequently used recipe with %v consumptions.\n",
			recipeName,
			recipeConsumptionCount,
		)
	}

	peakDayBytes, _ := json.Marshal(
		context.Highlights.PeakBatchDay,
	)

	var peakDay map[string]any

	json.Unmarshal(
		peakDayBytes,
		&peakDay,
	)

	peakDayName := peakDay["day_name"]
	peakConsumptionCount := peakDay["batch_count"]

	if peakDayName != nil {
		summary += fmt.Sprintf(
			"- %v recorded the highest operational activity with %v consumptions.\n",
			peakDayName,
			peakConsumptionCount,
		)
	}

	if len(context.TopRecipeVariance) > 0 {
		summary += fmt.Sprintf(
			"- %s recorded the highest recipe variance at %.1f%%.\n",
			context.TopRecipeVariance[0].RecipeName,
			context.TopRecipeVariance[0].VariancePercent,
		)
	}

	return summary
}

func BuildExecutiveSummaryExplanation(
	context ExecutiveAIContext,
) string {

	explanation := "I say that because the operational summary shows "

	explanation += fmt.Sprintf(
		"%d total consumptions across %d sites. ",
		context.KPIs.TotalBatches,
		context.KPIs.TotalBranches,
	)

	mostActiveSiteBytes, _ := json.Marshal(
		context.Highlights.MostActiveBranch,
	)

	var mostActiveSite map[string]any

	json.Unmarshal(
		mostActiveSiteBytes,
		&mostActiveSite,
	)

	siteName := mostActiveSite["branch_name"]
	consumptionCount := mostActiveSite["batch_count"]

	if siteName != nil {
		explanation += fmt.Sprintf(
			"%v currently accounts for %v consumptions and remains the busiest operational site. ",
			siteName,
			consumptionCount,
		)
	}

	explanation += fmt.Sprintf(
		"Daily planning also appears complete because there are %d finalized plans and %d drafts pending.",
		context.KPIs.FinalizedDailyPlans,
		context.KPIs.DraftDailyPlans,
	)

	return explanation
}

func BuildExecutiveSummaryRecommendation(
	context ExecutiveAIContext,
) string {

	recommendation := "Management should focus on operational concentration and recipe variance control. "

	mostActiveSiteBytes, _ := json.Marshal(
		context.Highlights.MostActiveBranch,
	)

	var mostActiveSite map[string]any

	json.Unmarshal(
		mostActiveSiteBytes,
		&mostActiveSite,
	)

	siteName := mostActiveSite["branch_name"]

	if siteName != nil {
		recommendation += fmt.Sprintf(
			"%v currently carries the highest operational activity and should be monitored for workload balance and operational resilience. ",
			siteName,
		)
	}

	if len(context.TopRecipeVariance) > 0 {
		recommendation += fmt.Sprintf(
			"%s should also be reviewed because it recorded the highest recipe variance at %.1f%%.",
			context.TopRecipeVariance[0].RecipeName,
			context.TopRecipeVariance[0].VariancePercent,
		)
	}

	return recommendation
}
