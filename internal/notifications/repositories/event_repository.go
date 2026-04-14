package repositories

import (
	"context"

	"newco-go-reporting-service/internal/notifications/dto"
	"newco-go-reporting-service/internal/notifications/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	GetUnprocessedEvents(ctx context.Context, limit int) ([]dto.ActivityEvent, error)
	MarkProcessed(ctx context.Context, eventID int64) error
}

type eventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) EventRepository {
	return &eventRepository{pool: pool}
}

func (r *eventRepository) GetUnprocessedEvents(ctx context.Context, limit int) ([]dto.ActivityEvent, error) {
	rows, err := r.pool.Query(ctx, queries.GetUnprocessedEvents, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []dto.ActivityEvent

	for rows.Next() {
		var e dto.ActivityEvent

		err := rows.Scan(
			&e.ID,
			&e.Action,
			&e.TargetType,
			&e.TargetID,
			&e.ActorStaffProfileID,
			&e.BranchID,
			&e.Message,
			&e.MetadataJSON,
			&e.CreatedAt,
			&e.ProcessedAt,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

func (r *eventRepository) MarkProcessed(ctx context.Context, eventID int64) error {
	_, err := r.pool.Exec(ctx, queries.MarkEventProcessed, eventID)
	return err
}
