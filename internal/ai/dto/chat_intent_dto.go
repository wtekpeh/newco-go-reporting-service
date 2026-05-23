package dto

type AIIntentClassification struct {
	Intent string `json:"intent"`

	ToolName string `json:"tool_name"`

	NeedsChart bool `json:"needs_chart"`

	ChartType string `json:"chart_type,omitempty"`

	Reason string `json:"reason"`
}
