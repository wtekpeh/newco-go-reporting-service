package dto

type ExecutiveSummaryRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	BranchID  *int64 `json:"branch_id,omitempty"`
}

type ExecutiveSummaryResponse struct {
	Message         string   `json:"message"`
	Summary         string   `json:"summary"`
	KeyInsights     []string `json:"key_insights"`
	Risks           []string `json:"risks"`
	Recommendations []string `json:"recommendations"`
	DataNotes       []string `json:"data_notes"`
}
