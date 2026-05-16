package dto

type BranchDashboardSummaryResponse struct {
	KPIs             BranchDashboardKPIs         `json:"kpis"`
	RecentDailyPlans []BranchRecentDailyPlanItem `json:"recent_daily_plans"`
}

type BranchDashboardKPIs struct {
	BranchID int64 `json:"branch_id"`

	TotalStaff int64 `json:"total_staff"`

	TotalBatches     int64 `json:"total_batches"`
	BatchesThisWeek  int64 `json:"batches_this_week"`
	BatchesThisMonth int64 `json:"batches_this_month"`

	TotalDailyPlans     int64 `json:"total_daily_plans"`
	FinalizedDailyPlans int64 `json:"finalized_daily_plans"`
	DraftDailyPlans     int64 `json:"draft_daily_plans"`
}

type BranchRecentDailyPlanItem struct {
	DailyPlanID int64  `json:"daily_plan_id"`
	UsedDate    string `json:"used_date"`
	Status      string `json:"status"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
}
