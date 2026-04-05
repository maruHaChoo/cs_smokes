package memory

import (
    "sync"

    "cs-smokes-bot/internal/domain"
)

type SessionRepository struct {
    mu       sync.RWMutex
    sessions map[int64]*domain.UserSession
}

func NewSessionRepository() *SessionRepository {
    return &SessionRepository{sessions: make(map[int64]*domain.UserSession)}
}

func (r *SessionRepository) GetOrCreate(userID, chatID int64) (*domain.UserSession, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if session, ok := r.sessions[userID]; ok {
        session.ChatID = chatID
        return cloneSession(session), nil
    }

    session := &domain.UserSession{UserID: userID, ChatID: chatID, CurrentScreen: domain.ScreenMain}
    r.sessions[userID] = cloneSession(session)
    return cloneSession(session), nil
}

func (r *SessionRepository) Save(session *domain.UserSession) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.sessions[session.UserID] = cloneSession(session)
    return nil
}

func cloneSession(s *domain.UserSession) *domain.UserSession {
    if s == nil { return nil }
    cp := *s
    if s.ContentMessageID != nil {
        value := *s.ContentMessageID
        cp.ContentMessageID = &value
    }
    return &cp
}
