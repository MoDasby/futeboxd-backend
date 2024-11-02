package repository_mocks

import (
	"errors"
	"time"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type MockSessionRepo struct {
	sessions map[string]domain.Session
}

func NewMockSessionRepo() *MockSessionRepo {
	return &MockSessionRepo{
		sessions: make(map[string]domain.Session),
	}
}

func (m *MockSessionRepo) FindOneByToken(token string) (*domain.Session, error) {
	session, exists := m.sessions[token]
	if !exists {
		return nil, errors.New("session not found")
	}
	return &session, nil
}

func (m *MockSessionRepo) Create(session *domain.Session) (*domain.Session, error) {
	session.CreatedAt = time.Now()
	m.sessions[session.Token] = *session
	return session, nil
}

func (m *MockSessionRepo) Update(session *domain.Session) error {
	if _, exists := m.sessions[session.Token]; !exists {
		return errors.New("session not found")
	}
	m.sessions[session.Token] = *session
	return nil
}

func (m *MockSessionRepo) Delete(id string) error {
	for k, v := range m.sessions {
		if v.ID == id {
			delete(m.sessions, k)

			return nil
		}
	}

	return errors.New("session not found")
}
