package repository

import (
	"context"
	"database/sql"

	"github.com/modasby/futeboxd-backend/core/pkg/pagination"

	"github.com/modasby/futeboxd-backend/core/internal/profile"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) profile.Repository {
	return &ProfileRepository{
		db: db,
	}
}

func (repo *ProfileRepository) FindOneByUsername(ctx context.Context, username, requesterID string) (*profile.Profile, error) {
	query := `
		SELECT 
			u.id, u.name, u.bio, u.username, u.favorite_team, u.profile_picture,
			(SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers_count,
            (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following_count,
			COUNT(r.id) as review_count,
            EXISTS(SELECT 1 FROM followers WHERE follower_id = $2 AND following_id = u.id) AS is_following
		FROM users u 
		LEFT JOIN reviews r ON r.user_id = u.id
		WHERE u.username = $1
		GROUP BY u.id
	`

	row := repo.db.QueryRowContext(ctx, query, username, requesterID)

	var profile profile.Profile
	var name, bio sql.NullString
	var favoriteTeam sql.NullInt64

	if err := row.Scan(
		&profile.UserID,
		&name,
		&bio,
		&profile.Username,
		&favoriteTeam,
		&profile.ProfilePicture,
		&profile.FollowersCount,
		&profile.FollowingCount,
		&profile.ReviewsCount,
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

	if bio.Valid {
		profile.Bio = bio.String
	}

	if name.Valid {
		profile.Name = name.String
	}

	if favoriteTeam.Valid {
		profile.FavoriteTeam = favoriteTeam.Int64
	}

	return &profile, nil
}

func (repo *ProfileRepository) Search(ctx context.Context, requesterID, term string, page *pagination.Page) ([]profile.Profile, error) {
	query := `
		SELECT 
			u.id, u.name, u.bio, u.username, u.favorite_team, u.profile_picture,
			(SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers_count,
            (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following_count,
			COUNT(r.id) as review_count,
            EXISTS(SELECT 1 FROM followers WHERE follower_id = $2 AND following_id = u.id) AS is_following
		FROM users u
		LEFT JOIN reviews r ON r.user_id = u.id
		WHERE u.username ILIKE $1
		GROUP BY u.id
		LIMIT $3
		OFFSET ($4 - 1) * $3
	`

	rows, err := repo.db.QueryContext(ctx, query, "%"+term+"%", requesterID, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	profiles := make([]profile.Profile, 0)

	for rows.Next() {
		var profile profile.Profile
		var name, bio sql.NullString
		var favoriteTeam sql.NullInt64

		if err := rows.Scan(
			&profile.UserID,
			&name,
			&bio,
			&profile.Username,
			&favoriteTeam,
			&profile.ProfilePicture,
			&profile.FollowersCount,
			&profile.FollowingCount,
			&profile.ReviewsCount,
			&profile.IsFollowing,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.NewHTTPErr(
					"Perfil não encontrado",
					404,
					"REPOSITORY:PROFILE:FIND_ONE_BY_USERNAME:NOT_FOUND",
				)
			}
			return nil, err
		}

		if bio.Valid {
			profile.Bio = bio.String
		}

		if name.Valid {
			profile.Name = name.String
		}

		if favoriteTeam.Valid {
			profile.FavoriteTeam = favoriteTeam.Int64
		}

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

func (repo *ProfileRepository) ListFollowers(ctx context.Context, requesterID, username string, page *pagination.Page) ([]profile.Profile, error) {
	query := `
		SELECT 
			u.id, u.name, u.bio, u.username, u.favorite_team, u.profile_picture,
			(SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers_count,
            (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following_count,
            CASE WHEN f.following_id = $1 THEN TRUE ELSE FALSE END AS is_following
		FROM followers f
		LEFT JOIN users u ON u.id = f.follower_id
		WHERE f.following_id = (select id from users where username = $2)
		GROUP BY f.following_id, u.id
		LIMIT $3
		OFFSET ($4 - 1) * $3
	`

	rows, err := repo.db.QueryContext(ctx, query, requesterID, username, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	profiles := make([]profile.Profile, 0)

	for rows.Next() {
		var profile profile.Profile
		var name, bio sql.NullString
		var favoriteTeam sql.NullInt64

		if err := rows.Scan(
			&profile.UserID,
			&name,
			&bio,
			&profile.Username,
			&favoriteTeam,
			&profile.ProfilePicture,
			&profile.FollowersCount,
			&profile.FollowingCount,
			&profile.IsFollowing,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.NewHTTPErr(
					"Perfil não encontrado",
					404,
					"REPOSITORY:PROFILE:FIND_ONE_BY_USERNAME:NOT_FOUND",
				)
			}
			return nil, err
		}

		if bio.Valid {
			profile.Bio = bio.String
		}

		if name.Valid {
			profile.Name = name.String
		}

		if favoriteTeam.Valid {
			profile.FavoriteTeam = favoriteTeam.Int64
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (repo *ProfileRepository) ListFollowing(ctx context.Context, requesterID, username string, page *pagination.Page) ([]profile.Profile, error) {
	query := `
		SELECT 
			u.id, u.name, u.bio, u.username, u.favorite_team, u.profile_picture,
			(SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers_count,
            (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following_count,
            CASE WHEN following_id = $1 THEN TRUE ELSE FALSE END AS is_following
		FROM followers f
		LEFT JOIN users u ON u.id = f.following_id
		WHERE f.follower_id = (select id from users where username = $2)
		GROUP BY f.following_id, u.id
		LIMIT $3
		OFFSET ($4 - 1) * $3
	`

	rows, err := repo.db.QueryContext(ctx, query, requesterID, username, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	profiles := make([]profile.Profile, 0)

	for rows.Next() {
		var profile profile.Profile
		var name, bio sql.NullString
		var favoriteTeam sql.NullInt64

		if err := rows.Scan(
			&profile.UserID,
			&name,
			&bio,
			&profile.Username,
			&favoriteTeam,
			&profile.ProfilePicture,
			&profile.FollowersCount,
			&profile.FollowingCount,
			&profile.IsFollowing,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.NewHTTPErr(
					"Perfil não encontrado",
					404,
					"REPOSITORY:PROFILE:FIND_ONE_BY_USERNAME:NOT_FOUND",
				)
			}
			return nil, err
		}

		if bio.Valid {
			profile.Bio = bio.String
		}

		if name.Valid {
			profile.Name = name.String
		}

		if favoriteTeam.Valid {
			profile.FavoriteTeam = favoriteTeam.Int64
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (repo *ProfileRepository) ListPopularProfiles(ctx context.Context, requesterID string, page *pagination.Page) ([]profile.Profile, error) {
	query := `
		SELECT 
			u.id, u.name, u.bio, u.username, u.favorite_team, u.profile_picture,
			(SELECT COUNT(*) FROM followers WHERE following_id = u.id) AS followers_count,
            (SELECT COUNT(*) FROM followers WHERE follower_id = u.id) AS following_count,
            CASE WHEN following_id = $1 THEN TRUE ELSE FALSE END AS is_following
		FROM followers f
		LEFT JOIN users u ON u.id = f.following_id
		WHERE f.created_at >= NOW() - INTERVAL '3 day'
		GROUP BY f.following_id, u.id
		ORDER BY following_count DESC
		LIMIT $2
		OFFSET ($3 - 1) * $2
	`

	rows, err := repo.db.QueryContext(ctx, query, requesterID, page.Size, page.Index)
	if err != nil {
		return nil, err
	}

	profiles := make([]profile.Profile, 0)

	for rows.Next() {
		var profile profile.Profile
		var name, bio sql.NullString
		var favoriteTeam sql.NullInt64

		if err := rows.Scan(
			&profile.UserID,
			&name,
			&bio,
			&profile.Username,
			&favoriteTeam,
			&profile.ProfilePicture,
			&profile.FollowersCount,
			&profile.FollowingCount,
			&profile.IsFollowing,
		); err != nil {
			return nil, err
		}

		if bio.Valid {
			profile.Bio = bio.String
		}

		if name.Valid {
			profile.Name = name.String
		}

		if favoriteTeam.Valid {
			profile.FavoriteTeam = favoriteTeam.Int64
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
