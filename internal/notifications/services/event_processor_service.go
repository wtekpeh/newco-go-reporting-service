package services

import (
	"context"
	"log"

	notificationsdto "newco-go-reporting-service/internal/notifications/dto"
	notificationsrepo "newco-go-reporting-service/internal/notifications/repositories"
)

type EventProcessorService interface {
	ProcessPendingEvents(ctx context.Context, limit int) error
}

type eventProcessorService struct {
	eventRepo        notificationsrepo.EventRepository
	notificationRepo notificationsrepo.NotificationRepository
}

func NewEventProcessorService(
	eventRepo notificationsrepo.EventRepository,
	notificationRepo notificationsrepo.NotificationRepository,
) EventProcessorService {
	return &eventProcessorService{
		eventRepo:        eventRepo,
		notificationRepo: notificationRepo,
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
		if err := s.notificationRepo.InsertNotification(
			ctx,
			event.ID,
			recipientID,
		); err != nil {
			return err
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
	case "cook_batch_created", "cook_batch_actuals_updated", "cook_batch_finalized":
		return true
	case "ingredient_scales_recalibrated":
		return false
	default:
		return false
	}
}
