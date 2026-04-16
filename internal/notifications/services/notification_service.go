package services

import (
	"context"

	"newco-go-reporting-service/internal/dto"
	notificationsdto "newco-go-reporting-service/internal/notifications/dto"
	notificationsrepo "newco-go-reporting-service/internal/notifications/repositories"
)

type NotificationService interface {
	ListMyNotifications(
		ctx context.Context,
		access *dto.AccessContext,
		limit int,
	) ([]notificationsdto.NotificationItem, error)

	GetMyUnreadCount(
		ctx context.Context,
		access *dto.AccessContext,
	) (int64, error)

	MarkMyNotificationAsRead(
		ctx context.Context,
		access *dto.AccessContext,
		notificationID int64,
	) error

	MarkAllMyNotificationsAsRead(
		ctx context.Context,
		access *dto.AccessContext,
	) error
}

type notificationService struct {
	repo notificationsrepo.NotificationRepository
}

func NewNotificationService(
	repo notificationsrepo.NotificationRepository,
) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) ListMyNotifications(
	ctx context.Context,
	access *dto.AccessContext,
	limit int,
) ([]notificationsdto.NotificationItem, error) {
	return s.repo.GetNotificationsByRecipientID(
		ctx,
		access.StaffProfileID,
		limit,
	)
}

func (s *notificationService) GetMyUnreadCount(
	ctx context.Context,
	access *dto.AccessContext,
) (int64, error) {
	return s.repo.GetUnreadNotificationCountByRecipientID(
		ctx,
		access.StaffProfileID,
	)
}

func (s *notificationService) MarkMyNotificationAsRead(
	ctx context.Context,
	access *dto.AccessContext,
	notificationID int64,
) error {
	return s.repo.MarkNotificationAsRead(
		ctx,
		notificationID,
		access.StaffProfileID,
	)
}

func (s *notificationService) MarkAllMyNotificationsAsRead(
	ctx context.Context,
	access *dto.AccessContext,
) error {
	return s.repo.MarkAllNotificationsAsReadByRecipientID(
		ctx,
		access.StaffProfileID,
	)
}
