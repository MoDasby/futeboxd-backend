package repository

import (
	"database/sql"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"

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
			u.id, u.name, u.bio, u.username, u.favorite_team,
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
		&profile.Name,
		&profile.Bio,
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

func (repo *ProfileRepository) Search(requesterID, term string, page *pagination.Page) ([]domain.Profile, error) {
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

	rows, err := repo.db.Query(query, "%"+term+"%", requesterID, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	profiles := make([]domain.Profile, 0)

	for rows.Next() {
		var profile domain.Profile

		if err := rows.Scan(
			&profile.UserID,
			&profile.Name,
			&profile.Bio,
			&profile.Username,
			&profile.FavoriteTeam,
			&profile.FollowersCount,
			&profile.FollowingCount,
			&profile.IsFollowing,
		); err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
