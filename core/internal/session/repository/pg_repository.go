package repository

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-backend/core/internal/session"
	errorsTypes "github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) session.Repository {
	return &sessionRepository{db: db}
}

func (repo *sessionRepository) FindOneByToken(ctx context.Context, token string) (*session.Session, error) {
	query := `
		SELECT s.id, s.token, s.user_id, s.expires_at, s.created_at
		FROM sessions s
		WHERE s.token = $1
	`

	row := repo.db.QueryRowContext(ctx, query, token)

	var ID, sessionToken, userID sql.NullString
	var expiresAt, createdAt time.Time

	if err := row.Scan(&ID, &sessionToken, &userID, &expiresAt, &createdAt); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, &errorsTypes.HTTPErr{
				Msg:        "Sessão inválida",
				Code:       http.StatusUnauthorized,
				Context:    "SESSION:REPOSITORY:FIND_ONE_BY_TOKEN:NOT_FOUND",
				StackTrace: errorsTypes.CaptureStackTrace(),
				ErrorCode:  utils.GetTraceIDFromCtx(ctx),
				Timestamp:  time.Now().UTC(),
				Original:   err,
			}
		}

		return nil, err
	}

	session := session.Session{
		ID:        ID.String,
		Token:     sessionToken.String,
		UserID:    userID.String,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}

	return &session, nil
}

func (repo *sessionRepository) Create(ctx context.Context, session *session.Session) (*session.Session, error) {
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
		INSERT INTO sessions (expires_at, token, user_id) 
		VALUES ($1, $2, $3) 
		RETURNING expires_at, created_at, token, user_id
	`

	row := tx.QueryRowContext(ctx, query, session.ExpiresAt, session.Token, session.UserID)

	if err := row.Scan(&session.ExpiresAt, &session.CreatedAt, &session.Token, &session.UserID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return session, nil
}

func (repo *sessionRepository) Update(ctx context.Context, session *session.Session) error {
	query := `
		UPDATE sessions
		SET expires_at = $2
		WHERE id = $1
	`

	if _, err := repo.db.ExecContext(ctx, query, session.ID, session.ExpiresAt); err != nil {
		return err
	}

	return nil
}

func (repo *sessionRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM sessions
		WHERE id = $1
	`

	if _, err := repo.db.Exec(query, id); err != nil {
		return err
	}

	return nil
}
