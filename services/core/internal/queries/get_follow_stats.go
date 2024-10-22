package queries

import (
	"database/sql"

	"github.com/modasby/futeboxd-api/services/core/internal/dto"
)

type FollowStatsQuery interface {
	GetFollowStats(userID, requesterID string) (*dto.FollowStatsDTO, error)
}

type followStatsQuery struct {
	db *sql.DB
}

func NewFollowStatsQuery(db *sql.DB) FollowStatsQuery {
	return &followStatsQuery{
		db: db,
	}
}

func (fs *followStatsQuery) GetFollowStats(userID, requesterID string) (*dto.FollowStatsDTO, error) {
	query := `
        SELECT 
            (SELECT COUNT(*) FROM followers WHERE following_id = $1) AS followers_count,
            (SELECT COUNT(*) FROM followers WHERE follower_id = $1) AS following_count,
            EXISTS(SELECT 1 FROM followers WHERE follower_id = $2 AND following_id = $1) AS is_following
    `

	row := fs.db.QueryRow(query, userID, requesterID)

	var stats dto.FollowStatsDTO
	err := row.Scan(&stats.FollowersCount, &stats.FollowingCount, &stats.Following)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
