package repository

import (
	"database/sql"
	"errors"
	"time"

	errorsTypes "github.com/modasby/futeboxd-api/pkg/errors"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) domain.SessionRepository {
	return &SessionRepository{db: db}
}

func (repo *SessionRepository) FindOneByToken(token string) (*domain.Session, error) {
	query := `
		SELECT s.id, s.token, s.user_id, s.expires_at, s.created_at
		FROM sessions s
		WHERE s.token = $1
	`

	row := repo.db.QueryRow(query, token)

	var ID, sessionToken, userID sql.NullString
	var expiresAt, createdAt time.Time

	if err := row.Scan(&ID, &sessionToken, &userID, &expiresAt, &createdAt); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorsTypes.NewHTTPErr(
				"sessão inválida",
				401,
				"REPOSITORY:SESSION:FIND_ONE_BY_TOKEN:NOT_FOUND",
			)
		}

		return nil, err
	}

	session := domain.Session{
		ID:        ID.String,
		Token:     sessionToken.String,
		UserID:    userID.String,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}

	return &session, nil
}

func (repo *SessionRepository) Create(session *domain.Session) (*domain.Session, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		WITH deleted AS (
			DELETE FROM sessions WHERE user_id = $3
		)

		INSERT INTO sessions (expires_at, token, user_id) 
		VALUES ($1, $2, $3) 
		RETURNING expires_at, created_at, token, user_id
	`

	row := tx.QueryRow(query, session.ExpiresAt, session.Token, session.UserID)

	if err := row.Scan(&session.ExpiresAt, &session.CreatedAt, &session.Token, &session.UserID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return session, nil
}

func (repo *SessionRepository) Update(session *domain.Session) error {
	query := `
		UPDATE sessions
		SET expires_at = $2
		WHERE id = $1
	`

	if _, err := repo.db.Exec(query, session.ID, session.ExpiresAt); err != nil {
		return err
	}

	return nil
}

func (repo *SessionRepository) Delete(id string) error {
	query := `
		DELETE FROM sessions
		WHERE id = $1
	`

	if _, err := repo.db.Exec(query, id); err != nil {
		return err
	}

	return nil
}
