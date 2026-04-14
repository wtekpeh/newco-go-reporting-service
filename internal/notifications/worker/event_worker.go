package worker

import (
	"context"
	"log"
	"time"

	notificationsservice "newco-go-reporting-service/internal/notifications/services"
)

type EventWorker struct {
	service  notificationsservice.EventProcessorService
	interval time.Duration
	limit    int
}

func NewEventWorker(
	service notificationsservice.EventProcessorService,
	interval time.Duration,
	limit int,
) *EventWorker {
	return &EventWorker{
		service:  service,
		interval: interval,
		limit:    limit,
	}
}

func (w *EventWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf(
		"[notifications] event worker started interval=%s limit=%d",
		w.interval.String(),
		w.limit,
	)

	for {
		select {
		case <-ctx.Done():
			log.Println("[notifications] event worker stopped")
			return

		case <-ticker.C:
			if err := w.service.ProcessPendingEvents(ctx, w.limit); err != nil {
				log.Printf("[notifications] worker error: %v", err)
			}
		}
	}
}
