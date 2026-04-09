package dto

type ExecutiveSummaryResponse struct {
	Message          string `json:"message"`
	TotalUsers       int    `json:"total_users"`
	TotalActiveUsers int    `json:"total_active_users"`
}
