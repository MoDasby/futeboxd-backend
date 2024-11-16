package repository_mocks

import (
	"errors"
	"time"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	errorsTypes "github.com/modasby/futeboxd-api/services/core/internal/errors"
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
		return nil, errorsTypes.NewHTTPErr(
			"sessão inválida",
			401,
			"REPOSITORY:SESSION:FIND_ONE_BY_TOKEN:NOT_FOUND",
		)
	}
	return &session, nil
}

func (m *MockSessionRepo) Create(session *domain.Session) (*domain.Session, error) {
	session.CreatedAt = time.Now()
	m.sessions[session.Token] = *session
	return session, nil
}

func (m *MockSessionRepo) Update(session *domain.Session) error {
	if _, exists := m.sessions[session.Token]; exists {
		m.sessions[session.Token] = *session
	}

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
