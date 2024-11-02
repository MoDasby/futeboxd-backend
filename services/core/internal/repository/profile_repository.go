package repository

import (
	"database/sql"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/errors"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) domain.ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (repo *ProfileRepository) FindOneByUsername(username, requesterID string) (*domain.Profile, error) {
	query := `
		SELECT 
			u.id, u.username, u.favorite_team,
			(SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers_count,
            (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following_count,
            EXISTS(SELECT 1 FROM followers WHERE follower_id = $2 AND following_id = u.id) AS is_following
		FROM users u 
		WHERE u.username = $1
	`

	row := repo.db.QueryRow(query, username, requesterID)

	var profile domain.Profile

	if err := row.Scan(
		&profile.UserID,
		&profile.Username,
		&profile.FavoriteTeam,
		&profile.FollowersCount,
		&profile.FollowingCount,
		&profile.IsFollowing,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewHTTPErr(
				"perfil não encontrado",
				404,
				"REPOSITORY:PROFILE:FIND_ONE_BY_USERNAME:NOT_FOUND",
			)
		}
		return nil, err
	}

	return &profile, nil
}
