package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/modasby/futeboxd-backend/core/internal/user"
	errorsTypes "github.com/modasby/futeboxd-backend/core/pkg/errors"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) FindBatchByID(ctx context.Context, ids []string) ([]user.User, error) {
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)) // slice de argumentos

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	placeholderStr := strings.Join(placeholders, ", ")

	query := fmt.Sprintf(`
		SELECT u.id::text, u.name, u.username, u.email, u.profile_picture
		FROM users u 
		WHERE u.id IN (%s)
	`, placeholderStr)

	rows, err := r.db.QueryContext(ctx, query, args...) // Use args com o operador variadic
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]user.User, 0)
	for rows.Next() {
		var user user.User
		var name sql.NullString

		if err := rows.Scan(&user.ID, &name, &user.Username, &user.Email, &user.ProfilePicture); err != nil {
			return nil, err
		}

		if name.Valid {
			user.Name = name.String
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *user.User) error {
	query := `
		UPDATE users
		SET name = $1, username = $2, email = $3, password = $4, favorite_team = $5, bio = $6, profile_picture = $7
		WHERE id = $8
	`

	if _, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.FavoriteTeamID,
		user.Bio,
		user.ProfilePicture,
		user.ID,
	); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindOneByIdOrUsername(ctx context.Context, identificator string) (*user.User, error) {
	query := `
		SELECT u.id::text, u.name, u.username, u.email, u.password, u.favorite_team, u.profile_picture, u.bio
		FROM users u
		WHERE LOWER(u.username) = LOWER($1) OR u.id::text = $1
	`

	rows := r.db.QueryRowContext(ctx, query, identificator)

	var user user.User
	var bio sql.NullString
	var name sql.NullString

	if err := rows.Scan(
		&user.ID,
		&name,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FavoriteTeamID,
		&user.ProfilePicture,
		&bio,
	); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorsTypes.NewHTTPErr(
				"usuário não encontrado",
				404,
				"REPOSITORY:USER:FIND_ONE_BY_ID_OR_USERNAME:NOT_FOUND",
			)
		}

		return nil, err
	}

	if bio.Valid {
		user.Bio = bio.String
	}

	if name.Valid {
		user.Name = name.String
	}

	return &user, nil
}

func (r *userRepository) FindOneByCredential(ctx context.Context, credential string) (*user.User, error) {
	query := `
		SELECT u.id::text, u.name, u.username, u.email, u.password, u.profile_picture
		FROM users u
		WHERE LOWER(u.email) = LOWER($1) OR LOWER(u.username) = LOWER($1)
	`
	rows := r.db.QueryRowContext(ctx, query, credential)

	var user user.User
	var name sql.NullString

	if err := rows.Scan(
		&user.ID,
		&name,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.ProfilePicture,
	); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errorsTypes.NewHTTPErr(
				"usuário não encontrado",
				404,
				"REPOSITORY:USER:FIND_ONE_BY_CREDENTIAL:NOT_FOUND",
			)
		}

		return nil, err
	}

	if name.Valid {
		user.Name = name.String
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *user.User) error {
	query := `
		INSERT INTO users(name, username, email, password, favorite_team, profile_picture)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx, query,
		user.Name, user.Username, user.Email, user.Password, user.FavoriteTeamID, user.ProfilePicture,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Exists(ctx context.Context, username, email string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1) OR LOWER(email) = LOWER($2))
	`

	row := r.db.QueryRowContext(ctx, query, username, email)

	var exists sql.NullBool

	if err := row.Scan(&exists); err != nil {

		return false, err
	}

	return exists.Bool, nil
}

func (repo *userRepository) SaveRecoverToken(ctx context.Context, recover *user.Recover) error {
	query := `
		INSERT INTO recover_password_tokens(user_id, token)
		VALUES ($1, $2)
	`

	_, err := repo.db.ExecContext(ctx, query, recover.UserID, recover.Token)

	return err
}

func (repo *userRepository) CheckRecoverToken(ctx context.Context, token string) (*user.Recover, error) {
	query := `
		SELECT user_id, token, expires_at
		FROM recover_password_tokens
		WHERE token = $1
	`

	row := repo.db.QueryRowContext(ctx, query, token)

	var recover user.Recover

	if err := row.Scan(&recover.UserID, &recover.Token, &recover.ExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsTypes.NewHTTPErr(
				"token inválido",
				404,
				"REPOSITORY:USER:CHECK_RECOVER_TOKEN:TOKEN_NOT_FOUND",
			)
		}
		return nil, err
	}

	return &recover, nil
}
