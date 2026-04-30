package dto

type ExecutiveSummaryResponse struct {
	Message    string                 `json:"message"`
	Kpis       ExecutiveKpisDTO       `json:"kpis"`
	Highlights ExecutiveHighlightsDTO `json:"highlights"`
}

type BranchActivityDTO struct {
	BranchName string `json:"branch_name"`
	BatchCount int    `json:"batch_count"`
}

type BranchStaffDTO struct {
	BranchName string `json:"branch_name"`
	StaffCount int    `json:"staff_count"`
}

type RecipeUsageDTO struct {
	RecipeName string `json:"recipe_name"`
	BatchCount int    `json:"batch_count"`
}

type AverageMetricDTO struct {
	Value float64 `json:"value"`
}

type PeakBatchDayDTO struct {
	DayName    string `json:"day_name"`
	BatchCount int    `json:"batch_count"`
}

type ExecutiveHighlightsDTO struct {
	MostActiveBranch        *BranchActivityDTO `json:"most_active_branch,omitempty"`
	LargestBranch           *BranchStaffDTO    `json:"largest_branch,omitempty"`
	MostUsedRecipe          *RecipeUsageDTO    `json:"most_used_recipe,omitempty"`
	AverageBatchesPerBranch *AverageMetricDTO  `json:"average_batches_per_branch,omitempty"`
	PeakBatchDay            *PeakBatchDayDTO   `json:"peak_batch_day,omitempty"`
}

type ExecutiveKpisDTO struct {
	TotalUsers       int `json:"total_users"`
	TotalActiveUsers int `json:"total_active_users"`
	TotalBranches    int `json:"total_branches"`
	TotalBatches     int `json:"total_batches"`
	BatchesThisWeek  int `json:"batches_this_week"`
	BatchesThisMonth int `json:"batches_this_month"`
}

type RecentBatchItemDTO struct {
	BatchID     int    `json:"batch_id"`
	RecipeName  string `json:"recipe_name"`
	BranchName  string `json:"branch_name"`
	CreatedBy   string `json:"created_by"`
	NPeople     int    `json:"n_people"`
	Status      string `json:"status"`
	ProteinType string `json:"protein_type"`
	CreatedAt   string `json:"created_at"`
}

type RecentBatchesResponse struct {
	Message string               `json:"message"`
	Items   []RecentBatchItemDTO `json:"items"`
}

type BatchTrendPointDTO struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

type BatchTrendsResponse struct {
	Message string               `json:"message"`
	Series  []BatchTrendPointDTO `json:"series"`
}

type BranchSummaryItemDTO struct {
	BranchID     int    `json:"branch_id"`
	BranchName   string `json:"branch_name"`
	StaffCount   int    `json:"staff_count"`
	TotalBatches int    `json:"total_batches"`
}

type BranchSummaryResponse struct {
	Message string                 `json:"message"`
	Items   []BranchSummaryItemDTO `json:"items"`
}

type RoleDistributionItemDTO struct {
	Role  string `json:"role"`
	Count int    `json:"count"`
}

type RoleDistributionResponse struct {
	Message string                    `json:"message"`
	Items   []RoleDistributionItemDTO `json:"items"`
}

type UserGrowthPointDTO struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

type UserGrowthResponse struct {
	Message string               `json:"message"`
	Series  []UserGrowthPointDTO `json:"series"`
}

type ReportFiltersDTO struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	BranchID  string `json:"branch_id"`
	GroupBy   string `json:"group_by"`
}

type GlobalRoleItemDTO struct {
	Role  string `json:"role"`
	Count int    `json:"count"`
}

type ActiveStatusItemDTO struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type StaffSummaryResponse struct {
	Message        string                    `json:"message"`
	UserTrends     []UserGrowthPointDTO      `json:"user_trends"`
	GlobalRoles    []GlobalRoleItemDTO       `json:"global_roles"`
	BranchRoles    []RoleDistributionItemDTO `json:"branch_roles"`
	ActiveStatuses []ActiveStatusItemDTO     `json:"active_statuses"`
}

type BatchStatusItemDTO struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type BatchSummaryResponse struct {
	Message      string               `json:"message"`
	TotalBatches int                  `json:"total_batches"`
	StatusCounts []BatchStatusItemDTO `json:"status_counts"`
}

type BranchTrendPointDTO struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

type BranchTrendSeriesDTO struct {
	BranchID   int                   `json:"branch_id"`
	BranchName string                `json:"branch_name"`
	Points     []BranchTrendPointDTO `json:"points"`
}

type BranchTrendsResponse struct {
	Message string                 `json:"message"`
	Series  []BranchTrendSeriesDTO `json:"series"`
}

type BranchItemDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type IngredientCategoryDailyDTO struct {
	CategoryID       *int64  `json:"category_id"`
	CategoryName     string  `json:"category_name"`
	Unit             string  `json:"unit"`
	TotalFinalValue  float64 `json:"total_final_value"`
	TotalActualValue float64 `json:"total_actual_value"`
}
