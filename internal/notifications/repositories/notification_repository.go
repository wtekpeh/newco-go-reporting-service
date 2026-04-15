package repositories

import (
	"context"
	"fmt"

	notificationsdto "newco-go-reporting-service/internal/notifications/dto"
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
	) (int64, error)

	GetNotificationsByRecipientID(
		ctx context.Context,
		recipientStaffProfileID int64,
		limit int,
	) ([]notificationsdto.NotificationItem, error)

	GetUnreadNotificationCountByRecipientID(
		ctx context.Context,
		recipientStaffProfileID int64,
	) (int64, error)

	MarkNotificationAsRead(
		ctx context.Context,
		notificationID int64,
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
) (int64, error) {
	var notificationID int64

	err := r.pool.QueryRow(
		ctx,
		queries.InsertNotification,
		eventID,
		recipientStaffProfileID,
	).Scan(&notificationID)
	if err != nil {
		return 0, err
	}

	return notificationID, nil
}
func (r *notificationRepository) GetNotificationsByRecipientID(
	ctx context.Context,
	recipientStaffProfileID int64,
	limit int,
) ([]notificationsdto.NotificationItem, error) {
	rows, err := r.pool.Query(
		ctx,
		queries.GetNotificationsByRecipientID,
		recipientStaffProfileID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]notificationsdto.NotificationItem, 0)

	for rows.Next() {
		var item notificationsdto.NotificationItem

		err := rows.Scan(
			&item.ID,
			&item.EventID,
			&item.Action,
			&item.TargetType,
			&item.TargetID,
			&item.Message,
			&item.IsRead,
			&item.ReadAt,
			&item.CreatedAt,
			&item.ActorStaffProfileID,
			&item.ActorFullName,
			&item.ActorEmail,
			&item.BranchID,
			&item.BranchName,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *notificationRepository) GetUnreadNotificationCountByRecipientID(
	ctx context.Context,
	recipientStaffProfileID int64,
) (int64, error) {
	var count int64

	err := r.pool.QueryRow(
		ctx,
		queries.GetUnreadNotificationCountByRecipientID,
		recipientStaffProfileID,
	).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *notificationRepository) MarkNotificationAsRead(
	ctx context.Context,
	notificationID int64,
	recipientStaffProfileID int64,
) error {
	tag, err := r.pool.Exec(
		ctx,
		queries.MarkNotificationAsRead,
		notificationID,
		recipientStaffProfileID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notification not found for this recipient")
	}

	return nil
}
