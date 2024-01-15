package database

import "sync"

type SessionManager struct {
	sessions map[string]string
	mutex    sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: map[string]string{},
	}
}

func (s *SessionManager) CreateSession(sessionId, sessionUser string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions[sessionId] = sessionUser
}

func (s *SessionManager) GetSession(sessionId string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.sessions[sessionId]
}

func (s *SessionManager) DeleteSession(sessionId string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, sessionId)
}
