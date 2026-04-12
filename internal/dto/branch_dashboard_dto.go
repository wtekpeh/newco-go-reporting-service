package dto

type BranchDashboardSummaryResponse struct {
	KPIs BranchDashboardKPIs `json:"kpis"`
}

type BranchDashboardKPIs struct {
	BranchID         int64 `json:"branch_id"`
	TotalStaff       int64 `json:"total_staff"`
	TotalBatches     int64 `json:"total_batches"`
	BatchesThisWeek  int64 `json:"batches_this_week"`
	BatchesThisMonth int64 `json:"batches_this_month"`
}
