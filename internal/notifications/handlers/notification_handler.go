package handlers

import (
	"log"
	"strconv"

	"newco-go-reporting-service/internal/dto"
	notificationsdto "newco-go-reporting-service/internal/notifications/dto"
	notificationsservice "newco-go-reporting-service/internal/notifications/services"

	"github.com/gofiber/fiber/v2"
)

const accessContextKey = "access_context"

type NotificationHandler struct {
	service notificationsservice.NotificationService
}

func NewNotificationHandler(
	service notificationsservice.NotificationService,
) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) ListMyNotifications(c *fiber.Ctx) error {
	value := c.Locals(accessContextKey)
	access, ok := value.(*dto.AccessContext)
	if !ok || access == nil {
		return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
			Message: "access context missing",
		})
	}

	log.Printf(
		"[notifications] list request staff_profile_id=%d email=%s role=%s",
		access.StaffProfileID,
		access.Email,
		access.GlobalRole,
	)

	limit := 20
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	items, err := h.service.ListMyNotifications(c.UserContext(), access, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch notifications",
			Error:   err.Error(),
		})
	}

	return c.JSON(notificationsdto.NotificationListResponse{
		Message: "notifications fetched successfully",
		Items:   items,
	})
}

func (h *NotificationHandler) GetMyUnreadCount(c *fiber.Ctx) error {
	value := c.Locals(accessContextKey)
	access, ok := value.(*dto.AccessContext)
	if !ok || access == nil {
		return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
			Message: "access context missing",
		})
	}

	log.Printf(
		"[notifications] unread-count request staff_profile_id=%d email=%s role=%s",
		access.StaffProfileID,
		access.Email,
		access.GlobalRole,
	)

	count, err := h.service.GetMyUnreadCount(c.UserContext(), access)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to fetch unread count",
			Error:   err.Error(),
		})
	}

	return c.JSON(notificationsdto.UnreadCountResponse{
		Message:     "unread count fetched successfully",
		UnreadCount: count,
	})
}

func (h *NotificationHandler) MarkMyNotificationAsRead(c *fiber.Ctx) error {
	value := c.Locals(accessContextKey)
	access, ok := value.(*dto.AccessContext)
	if !ok || access == nil {
		return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
			Message: "access context missing",
		})
	}

	notificationID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || notificationID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "invalid notification id",
		})
	}

	if err := h.service.MarkMyNotificationAsRead(
		c.UserContext(),
		access,
		notificationID,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "failed to mark notification as read",
			Error:   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "notification marked as read",
	})
}
