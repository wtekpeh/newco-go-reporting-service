package dto

type AIIntentClassification struct {
	Intent string `json:"intent"`

	ToolName string `json:"tool_name"`

	ReasoningMode string `json:"reasoning_mode"`

	NeedsChart bool `json:"needs_chart"`

	ChartType string `json:"chart_type,omitempty"`

	Reason string `json:"reason"`

	ConfidenceScore float64 `json:"confidence_score"`

	ClassificationSource string `json:"classification_source"`
}
