package repository

import (
	"context"
	"database/sql"

	"github.com/modasby/futeboxd-backend/core/pkg/pagination"

	"github.com/modasby/futeboxd-backend/core/internal/domain"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) domain.ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (repo *ProfileRepository) FindOneByUsername(ctx context.Context, username, requesterID string) (*domain.Profile, error) {
	query := `
		SELECT 
			u.id, u.name, u.bio, u.username, u.favorite_team,
			(SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers_count,
            (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following_count,
            EXISTS(SELECT 1 FROM followers WHERE follower_id = $2 AND following_id = u.id) AS is_following
		FROM users u 
		WHERE u.username = $1
	`

	row := repo.db.QueryRowContext(ctx, query, username, requesterID)

	var profile domain.Profile
	var bio sql.NullString

	if err := row.Scan(
		&profile.UserID,
		&profile.Name,
		&bio,
		&profile.Username,
		&profile.FavoriteTeam,
		&profile.FollowersCount,
		&profile.FollowingCount,
		&profile.IsFollowing,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr(
				"perfil "+username+" não encontrado",
				404,
				"REPOSITORY:PROFILE:FIND_ONE_BY_USERNAME:NOT_FOUND",
			)
		}
		return nil, err
	}

	profile.Bio = bio.String

	return &profile, nil
}

func (repo *ProfileRepository) Search(ctx context.Context, requesterID, term string, page *pagination.Page) ([]domain.Profile, error) {
	query := `
		SELECT 
			u.id, u.name, u.bio, u.username, u.favorite_team,
			(SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers_count,
            (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following_count,
            EXISTS(SELECT 1 FROM followers WHERE follower_id = $2 AND following_id = u.id) AS is_following
		FROM users u
		WHERE u.username ILIKE $1
		LIMIT $3
		OFFSET ($4 - 1) * $3
	`

	rows, err := repo.db.QueryContext(ctx, query, "%"+term+"%", requesterID, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	profiles := make([]domain.Profile, 0)

	for rows.Next() {
		var profile domain.Profile
		var bio sql.NullString

		if err := rows.Scan(
			&profile.UserID,
			&profile.Name,
			&bio,
			&profile.Username,
			&profile.FavoriteTeam,
			&profile.FollowersCount,
			&profile.FollowingCount,
			&profile.IsFollowing,
		); err != nil {
			return nil, err
		}

		profile.Bio = bio.String

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (r *ProfileRepository) Follow(ctx context.Context, followerID, followingID string) error {
	query := `
		INSERT INTO followers (follower_id, following_id)
		VALUES ($1, $2);
	`

	if _, err := r.db.ExecContext(ctx, query, followerID, followingID); err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) Unfollow(ctx context.Context, followerID, followingID string) error {
	query := `
		DELETE FROM followers
		WHERE follower_id = $1 AND following_id = $2
	`

	if _, err := r.db.ExecContext(ctx, query, followerID, followingID); err != nil {
		return err
	}

	return nil
}

func (r *ProfileRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM followers
			WHERE follower_id = $1 AND following_id = $2
		)
	`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, followerID, followingID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
