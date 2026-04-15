package dto

type LiveNotificationMessage struct {
	Type           string `json:"type"`
	NotificationID int64  `json:"notification_id"`
	EventID        int64  `json:"event_id"`
	Action         string `json:"action"`
	TargetType     string `json:"target_type"`
	TargetID       *int64 `json:"target_id,omitempty"`
	Message        string `json:"message"`
	BranchID       *int64 `json:"branch_id,omitempty"`
	RecipientID    int64  `json:"recipient_id"`
	CreatedAt      string `json:"created_at"`
}
