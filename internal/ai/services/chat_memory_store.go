package services

import (
	"sync"

	"newco-go-reporting-service/internal/ai/dto"
)

type ChatMemoryStore struct {
	mu       sync.RWMutex
	sessions map[string][]dto.AIConversationTurn
}

func NewChatMemoryStore() *ChatMemoryStore {
	return &ChatMemoryStore{
		sessions: make(map[string][]dto.AIConversationTurn),
	}
}

func (s *ChatMemoryStore) AddTurn(turn dto.AIConversationTurn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[turn.SessionID] = append(
		s.sessions[turn.SessionID],
		turn,
	)

	if len(s.sessions[turn.SessionID]) > 10 {
		s.sessions[turn.SessionID] =
			s.sessions[turn.SessionID][len(s.sessions[turn.SessionID])-10:]
	}
}

func (s *ChatMemoryStore) RecentTurns(sessionID string) []dto.AIConversationTurn {
	s.mu.RLock()
	defer s.mu.RUnlock()

	turns := s.sessions[sessionID]

	copyOfTurns := make([]dto.AIConversationTurn, len(turns))
	copy(copyOfTurns, turns)

	return copyOfTurns
}
