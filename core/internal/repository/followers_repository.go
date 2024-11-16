package repository

import (
	"database/sql"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

type FollowersRepository struct {
	db *sql.DB
}

func NewFollowersRepository(db *sql.DB) domain.FollowersRepository {
	return &FollowersRepository{db: db}
}

func (r *FollowersRepository) Follow(followerID, followingID string) error {
	query := `
		INSERT INTO followers (follower_id, following_id)
		VALUES ($1, $2);
	`

	if _, err := r.db.Exec(query, followerID, followingID); err != nil {
		return err
	}

	return nil
}

func (r *FollowersRepository) Unfollow(followerID, followingID string) error {
	query := `
		DELETE FROM followers
		WHERE follower_id = $1 AND following_id = $2
	`

	if _, err := r.db.Exec(query, followerID, followingID); err != nil {
		return err
	}

	return nil
}

func (r *FollowersRepository) IsFollowing(followerID, followingID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM followers
			WHERE follower_id = $1 AND following_id = $2
		)
	`

	var exists bool
	if err := r.db.QueryRow(query, followerID, followingID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
