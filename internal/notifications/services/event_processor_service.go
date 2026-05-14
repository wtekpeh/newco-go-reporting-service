package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	notificationsdto "newco-go-reporting-service/internal/notifications/dto"
	notificationsrealtime "newco-go-reporting-service/internal/notifications/realtime"
	notificationsrepo "newco-go-reporting-service/internal/notifications/repositories"
)

type EventProcessorService interface {
	ProcessPendingEvents(ctx context.Context, limit int) error
}

type eventProcessorService struct {
	eventRepo        notificationsrepo.EventRepository
	notificationRepo notificationsrepo.NotificationRepository
	hub              *notificationsrealtime.Hub
}

func NewEventProcessorService(
	eventRepo notificationsrepo.EventRepository,
	notificationRepo notificationsrepo.NotificationRepository,
	hub *notificationsrealtime.Hub,
) EventProcessorService {
	return &eventProcessorService{
		eventRepo:        eventRepo,
		notificationRepo: notificationRepo,
		hub:              hub,
	}
}

func (s *eventProcessorService) ProcessPendingEvents(
	ctx context.Context,
	limit int,
) error {
	events, err := s.eventRepo.GetUnprocessedEvents(ctx, limit)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := s.processSingleEvent(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func (s *eventProcessorService) processSingleEvent(
	ctx context.Context,
	event notificationsdto.ActivityEvent,
) error {
	log.Printf(
		"[notifications] processing event id=%d action=%s target_type=%s",
		event.ID,
		event.Action,
		event.TargetType,
	)

	recipientIDs, err := s.resolveRecipientIDs(ctx, event)
	if err != nil {
		return err
	}

	for _, recipientID := range recipientIDs {
		notificationID, err := s.notificationRepo.InsertNotification(
			ctx,
			event.ID,
			recipientID,
		)
		if err != nil {
			return err
		}

		if err := s.pushLiveNotification(event, recipientID, notificationID); err != nil {
			log.Printf(
				"[notifications] live push failed event_id=%d recipient_id=%d err=%v",
				event.ID,
				recipientID,
				err,
			)
		}
	}

	if err := s.eventRepo.MarkProcessed(ctx, event.ID); err != nil {
		return err
	}

	log.Printf(
		"[notifications] event id=%d processed with %d recipients",
		event.ID,
		len(recipientIDs),
	)

	return nil
}

func (s *eventProcessorService) resolveRecipientIDs(
	ctx context.Context,
	event notificationsdto.ActivityEvent,
) ([]int64, error) {
	recipientSet := make(map[int64]struct{})

	executiveIDs, err := s.notificationRepo.GetExecutiveRecipientIDs(ctx)
	if err != nil {
		return nil, err
	}

	for _, id := range executiveIDs {
		recipientSet[id] = struct{}{}
	}

	if s.shouldIncludeBranchManagers(event) && event.BranchID != nil {
		branchManagerIDs, err := s.notificationRepo.GetBranchManagerRecipientIDsByBranchID(
			ctx,
			*event.BranchID,
		)
		if err != nil {
			return nil, err
		}

		for _, id := range branchManagerIDs {
			recipientSet[id] = struct{}{}
		}
	}

	recipientIDs := make([]int64, 0, len(recipientSet))
	for id := range recipientSet {
		recipientIDs = append(recipientIDs, id)
	}

	return recipientIDs, nil
}

func (s *eventProcessorService) shouldIncludeBranchManagers(
	event notificationsdto.ActivityEvent,
) bool {
	switch event.Action {
	case "cook_batch_created",
		"cook_batch_actuals_updated",
		"cook_batch_finalized",

		"daily_plan_created",
		"daily_plan_actuals_updated",
		"daily_plan_finalized":
		return true
	case "ingredient_scales_recalibrated":
		return false
	default:
		return false
	}
}

func (s *eventProcessorService) pushLiveNotification(
	event notificationsdto.ActivityEvent,
	recipientID int64,
	notificationID int64,
) error {
	payload := notificationsdto.LiveNotificationMessage{
		Type:           "notification",
		NotificationID: notificationID,
		EventID:        event.ID,
		Action:         event.Action,
		TargetType:     event.TargetType,
		TargetID:       event.TargetID,
		Message:        event.Message,
		BranchID:       event.BranchID,
		RecipientID:    recipientID,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	s.hub.SendToStaffProfile(recipientID, data)
	return nil
}
