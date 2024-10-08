package repository

import (
	"database/sql"
	"errors"
	"time"

	errorsTypes "github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/users/internal/domain"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) domain.SessionRepository {
	return &SessionRepository{db: db}
}

func (repo *SessionRepository) FindOneByToken(token string) (*domain.Session, error) {
	query := `
		SELECT s.token, s.user_id, s.expires_at, s.created_at, u.username, u.email
		FROM sessions s
		INNER JOIN users u ON s.user_id = u.id
		WHERE s.token = $1
	`

	row := repo.db.QueryRow(query, token)

	var sessionToken, userID, username, email sql.NullString
	var expiresAt, createdAt time.Time

	if err := row.Scan(&sessionToken, &userID, &expiresAt, &createdAt, &username, &email); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorsTypes.NewErrUnauthorized("sessão inválida")
		}

		return nil, err
	}

	session := domain.NewSession(
		sessionToken.String,
		&domain.User{ID: userID.String, Username: username.String, Email: email.String},
		expiresAt,
	)
	session.CreatedAt = createdAt

	return session, nil
}

func (repo *SessionRepository) AddSession(session *domain.Session) (*domain.Session, error) {
	query := `
		INSERT INTO sessions (expires_at, token, user_id)
		VALUES ($1, $2, $3)
		RETURNING expires_at, created_at, token, user_id
	`

	row := repo.db.QueryRow(query, session.ExpiresAt, session.Token, session.Owner.ID)

	if err := row.Scan(&session.ExpiresAt, &session.CreatedAt, &session.Token, &session.Owner.ID); err != nil {
		return nil, err
	}

	return session, nil
}

func (repo *SessionRepository) DeleteSession(id string) error {
	return nil
}
