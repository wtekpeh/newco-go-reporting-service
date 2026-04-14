package repositories

import (
	"context"

	"newco-go-reporting-service/internal/notifications/queries"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository interface {
	GetExecutiveRecipientIDs(ctx context.Context) ([]int64, error)
	GetBranchManagerRecipientIDsByBranchID(
		ctx context.Context,
		branchID int64,
	) ([]int64, error)
	InsertNotification(
		ctx context.Context,
		eventID int64,
		recipientStaffProfileID int64,
	) error
}

type notificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) NotificationRepository {
	return &notificationRepository{pool: pool}
}

func (r *notificationRepository) GetExecutiveRecipientIDs(
	ctx context.Context,
) ([]int64, error) {
	rows, err := r.pool.Query(ctx, queries.GetExecutiveRecipientIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *notificationRepository) GetBranchManagerRecipientIDsByBranchID(
	ctx context.Context,
	branchID int64,
) ([]int64, error) {
	rows, err := r.pool.Query(
		ctx,
		queries.GetBranchManagerRecipientIDsByBranchID,
		branchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *notificationRepository) InsertNotification(
	ctx context.Context,
	eventID int64,
	recipientStaffProfileID int64,
) error {
	_, err := r.pool.Exec(
		ctx,
		queries.InsertNotification,
		eventID,
		recipientStaffProfileID,
	)
	return err
}
