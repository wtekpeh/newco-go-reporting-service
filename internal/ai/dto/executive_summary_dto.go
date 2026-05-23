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

type AIChatRequest struct {
	SessionID string `json:"session_id"`

	Message string `json:"message"`

	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`

	BranchID *int64 `json:"branch_id,omitempty"`
}

type AIChatResponse struct {
	Message string `json:"message"`

	SessionID string `json:"session_id"`

	Intent string `json:"intent"`

	UsedTools []string `json:"used_tools"`

	AssistantResponse string `json:"assistant_response"`

	ChartSuggestions []AIChartSuggestion `json:"chart_suggestions,omitempty"`

	ChartData map[string]interface{} `json:"chart_data,omitempty"`

	DataNotes []string `json:"data_notes,omitempty"`
}

type AIChartSuggestion struct {
	ChartType string `json:"chart_type"`

	Title string `json:"title"`

	Dataset string `json:"dataset"`

	XField string `json:"x_field"`

	YField string `json:"y_field"`
}
