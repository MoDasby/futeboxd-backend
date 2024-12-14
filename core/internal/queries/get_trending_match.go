package queries

import (
	"database/sql"

	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
)

type ListTrendingMatchQuery struct {
	db             *sql.DB
	footballClient football.Client
}

func NewListTrendingMatchQuery(
	db *sql.DB,
	footballClient football.Client,
) *ListTrendingMatchQuery {
	return &ListTrendingMatchQuery{
		db:             db,
		footballClient: footballClient,
	}
}

type RatingCount struct {
	Rating int `json:"rating"`
	Count  int `json:"count"`
}

type RateStats struct {
	RatingCount [5]RatingCount `json:"rating_count"`
	Avg         float32        `json:"avg"`
}

type TrendingMatchOutputDTO struct {
	football.Match
	RateStats *RateStats `json:"rate_stats"`
}

type TrendingMatchDTO struct {
	Match     int64
	RateStats *RateStats
}

func (tq *ListTrendingMatchQuery) Execute(page *pagination.Page) ([]TrendingMatchOutputDTO, error) {
	query := `
		SELECT 
		    r.match_id, 
		    ROUND(AVG(r.rate), 1) as average_rate,
		    SUM(CASE WHEN r.rate = 1 THEN 1 ELSE 0 END) AS count1,
  			SUM(CASE WHEN r.rate = 2 THEN 1 ELSE 0 END) AS count2,
			SUM(CASE WHEN r.rate = 3 THEN 1 ELSE 0 END) AS count3,
			SUM(CASE WHEN r.rate = 4 THEN 1 ELSE 0 END) AS count4,
			SUM(CASE WHEN r.rate = 5 THEN 1 ELSE 0 END) AS count5
		FROM reviews r
		WHERE r.created_at >= NOW() - INTERVAL '7 days'
		GROUP BY r.match_id
		ORDER BY (
    		LEAST(COUNT(r.match_id), 100)
		) + (
			ROUND(AVG(r.rate), 1)
		) DESC
		LIMIT $1
		OFFSET ($2 - 1) * $1
	`

	trendingMatches := make([]TrendingMatchDTO, 0)
	ids := make([]int64, 0)

	rows, err := tq.db.Query(query, page.Size, page.Index)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var matchID int64
		var avgRate float32
		var count1, count2, count3, count4, count5 int

		if err := rows.Scan(
			&matchID,
			&avgRate,
			&count1,
			&count2,
			&count3,
			&count4,
			&count5,
		); err != nil {
			return nil, err
		}

		trendingMatches = append(trendingMatches, TrendingMatchDTO{
			Match: matchID,
			RateStats: &RateStats{
				RatingCount: [5]RatingCount{
					{Count: count1, Rating: 1},
					{Count: count2, Rating: 2},
					{Count: count3, Rating: 3},
					{Count: count4, Rating: 4},
					{Count: count5, Rating: 5},
				},
				Avg: avgRate,
			},
		})
		ids = append(ids, matchID)
	}

	matches, err := tq.footballClient.GetMatchesMap(ids)
	if err != nil {
		return nil, err
	}

	rateMap := make(map[int64]*RateStats)

	for _, trendingMatch := range trendingMatches {
		rateMap[trendingMatch.Match] = trendingMatch.RateStats
	}

	output := make([]TrendingMatchOutputDTO, len(trendingMatches))

	for index, id := range ids {
		output[index] = TrendingMatchOutputDTO{
			Match:     matches[id],
			RateStats: rateMap[id],
		}
	}

	return output, nil
}
