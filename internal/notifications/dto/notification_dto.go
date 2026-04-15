package dto

import "time"

type NotificationItem struct {
	ID                  int64      `json:"id"`
	EventID             int64      `json:"event_id"`
	Action              string     `json:"action"`
	TargetType          string     `json:"target_type"`
	TargetID            *int64     `json:"target_id,omitempty"`
	Message             string     `json:"message"`
	IsRead              bool       `json:"is_read"`
	CreatedAt           time.Time  `json:"created_at"`
	ReadAt              *time.Time `json:"read_at,omitempty"`
	ActorStaffProfileID int64      `json:"actor_staff_profile_id"`
	ActorFullName       string     `json:"actor_full_name"`
	ActorEmail          string     `json:"actor_email"`
	BranchID            *int64     `json:"branch_id,omitempty"`
	BranchName          *string    `json:"branch_name,omitempty"`
}

type NotificationListResponse struct {
	Message string             `json:"message"`
	Items   []NotificationItem `json:"items"`
}

type UnreadCountResponse struct {
	Message     string `json:"message"`
	UnreadCount int64  `json:"unread_count"`
}
