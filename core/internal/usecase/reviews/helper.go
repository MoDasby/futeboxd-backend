package usecase

import (
	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
)

func GetMatchesByReview(footballClient football.Client, reviews []domain.Review) (map[int64]football.Match, error) {
	matchIDs := make([]int64, len(reviews))

	for i, review := range reviews {
		matchIDs[i] = review.MatchID
	}

	return footballClient.GetMatchesMap(matchIDs)
}
