package store

import (
	"visio/internal/database"
)

type Sessions struct {
	manager *database.SessionManager
}

func NewSessionsStore(manager *database.SessionManager) *Sessions {
	return &Sessions{
		manager: manager,
	}
}

func (s *Sessions) Create(id, sessionUser string) {
	s.manager.CreateSession(id, sessionUser)
}

func (s *Sessions) Get(id string) string {
	return s.manager.GetSession(id)
}

func (s *Sessions) Delete(id string) {
	s.manager.DeleteSession(id)
}
