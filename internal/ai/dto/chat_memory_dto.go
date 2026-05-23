package dto

import "time"

type AIConversationTurn struct {
	SessionID string `json:"session_id"`

	UserMessage string `json:"user_message"`

	Intent string `json:"intent"`

	ToolName string `json:"tool_name"`

	StartDate string `json:"start_date"`

	EndDate string `json:"end_date"`

	BranchID *int64 `json:"branch_id,omitempty"`

	AssistantResponse string `json:"assistant_response"`

	CreatedAt time.Time `json:"created_at"`
}
