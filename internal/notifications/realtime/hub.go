package realtime

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	mu          sync.RWMutex
	connections map[int64]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[int64]map[*websocket.Conn]struct{}),
	}
}

func (h *Hub) Register(staffProfileID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.connections[staffProfileID]; !exists {
		h.connections[staffProfileID] = make(map[*websocket.Conn]struct{})
	}

	h.connections[staffProfileID][conn] = struct{}{}

	log.Printf(
		"[notifications] websocket connected staff_profile_id=%d total_connections=%d",
		staffProfileID,
		len(h.connections[staffProfileID]),
	)
}

func (h *Hub) Unregister(staffProfileID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	userConnections, exists := h.connections[staffProfileID]
	if !exists {
		return
	}

	delete(userConnections, conn)

	if len(userConnections) == 0 {
		delete(h.connections, staffProfileID)
	}

	log.Printf(
		"[notifications] websocket disconnected staff_profile_id=%d remaining_connections=%d",
		staffProfileID,
		len(userConnections),
	)
}

func (h *Hub) SendToStaffProfile(staffProfileID int64, payload []byte) {
	h.mu.RLock()
	userConnections, exists := h.connections[staffProfileID]
	if !exists {
		h.mu.RUnlock()
		return
	}

	conns := make([]*websocket.Conn, 0, len(userConnections))
	for conn := range userConnections {
		conns = append(conns, conn)
	}
	h.mu.RUnlock()

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Printf(
				"[notifications] websocket write failed staff_profile_id=%d err=%v",
				staffProfileID,
				err,
			)
		}
	}
}
