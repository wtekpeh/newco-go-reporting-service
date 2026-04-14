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
	eventRepo notificationsrepo.EventRepository
}

func NewEventProcessorService(
	eventRepo notificationsrepo.EventRepository,
) EventProcessorService {
	return &eventProcessorService{
		eventRepo: eventRepo,
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

	if err := s.eventRepo.MarkProcessed(ctx, event.ID); err != nil {
		return err
	}

	log.Printf(
		"[notifications] marked event id=%d as processed",
		event.ID,
	)

	return nil
}
