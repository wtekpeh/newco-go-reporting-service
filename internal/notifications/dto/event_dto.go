package dto

import "time"

type ActivityEvent struct {
	ID         int64
	Action     string
	TargetType string
	TargetID   *int64

	ActorStaffProfileID int64
	BranchID            *int64

	Message string

	MetadataJSON map[string]interface{}

	CreatedAt   time.Time
	ProcessedAt *time.Time
}
