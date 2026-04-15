package handlers

import (
	"log"
	"strconv"

	"newco-go-reporting-service/internal/dto"
	notificationsrealtime "newco-go-reporting-service/internal/notifications/realtime"

	"github.com/gofiber/contrib/websocket"
)

const wsAccessContextKey = "access_context"

type NotificationWebSocketHandler struct {
	hub *notificationsrealtime.Hub
}

func NewNotificationWebSocketHandler(
	hub *notificationsrealtime.Hub,
) *NotificationWebSocketHandler {
	return &NotificationWebSocketHandler{
		hub: hub,
	}
}

func (h *NotificationWebSocketHandler) Handle(conn *websocket.Conn) {
	value := conn.Locals(wsAccessContextKey)
	access, ok := value.(*dto.AccessContext)
	if !ok || access == nil {
		log.Println("[notifications] websocket missing access context")
		_ = conn.Close()
		return
	}

	staffProfileID := access.StaffProfileID
	h.hub.Register(staffProfileID, conn)
	defer h.hub.Unregister(staffProfileID, conn)

	log.Printf(
		"[notifications] websocket session opened staff_profile_id=%d",
		staffProfileID,
	)

	welcomePayload := []byte(
		`{"type":"connected","staff_profile_id":` +
			strconv.FormatInt(staffProfileID, 10) +
			`}`,
	)

	if err := conn.WriteMessage(websocket.TextMessage, welcomePayload); err != nil {
		log.Printf(
			"[notifications] websocket welcome write failed staff_profile_id=%d err=%v",
			staffProfileID,
			err,
		)
		return
	}

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf(
				"[notifications] websocket read ended staff_profile_id=%d err=%v",
				staffProfileID,
				err,
			)
			return
		}

		if messageType == websocket.CloseMessage {
			log.Printf(
				"[notifications] websocket close received staff_profile_id=%d",
				staffProfileID,
			)
			return
		}
	}
}
