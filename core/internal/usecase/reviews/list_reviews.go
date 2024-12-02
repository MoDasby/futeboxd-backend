package usecase

import (
	"fmt"
	"strings"

	"github.com/modasby/futeboxd-api/services/core/internal/client/football"
	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/dto"
	"github.com/modasby/futeboxd-api/services/core/internal/pagination"
)

type ListReviewsUseCase struct {
	reviewRepository domain.ReviewRepository
	footballClient   football.Client
}

func NewListReviewsUseCase(
	reviewRepository domain.ReviewRepository,
	footballClient football.Client,
) *ListReviewsUseCase {
	return &ListReviewsUseCase{
		reviewRepository: reviewRepository,
		footballClient:   footballClient,
	}
}

func (uc ListReviewsUseCase) Execute(
	requesterID string,
	options ListReviewsOptions,
	page *pagination.Page,
) ([]dto.ReviewDTO, error) {

	whereClause, params := uc.buildWhereClause(options)

	reviews, err := uc.reviewRepository.ListAll(requesterID, whereClause, params, page)
	if err != nil {
		return nil, err
	}

	matches, err := GetMatchesByReview(uc.footballClient, reviews)
	if err != nil {
		return nil, err
	}

	output := make([]dto.ReviewDTO, 0)

	for _, review := range reviews {
		match := matches[review.MatchID]

		output = append(output, dto.ReviewDTO{
			ID: review.ID,
			Author: dto.UserDTO{
				ID:       review.Author.ID,
				Username: review.Author.Username,
			},
			Rate:          review.Rate,
			Description:   review.Description,
			Match:         &match,
			CommentsCount: review.CommentsCount,
			Likes:         review.Likes,
			IsLiked:       review.IsLiked,
			CreatedAt:     review.CreatedAt,
		})
	}

	return output, nil
}

type ListReviewsOptions struct {
	Username   string
	Team       string
	Match      string
	SearchTerm string
}

func (uc *ListReviewsUseCase) buildWhereClause(options ListReviewsOptions) (string, []any) {
	var whereClause strings.Builder
	var params []any
	counter := 1
	isFirstCondition := true

	appendCondition := func(condition, operator string, param any) {
		if !isFirstCondition {
			whereClause.WriteString(" " + operator + " ")
		}
		whereClause.WriteString(condition)
		params = append(params, param)
		isFirstCondition = false
		counter++
	}

	if options.Match != "" {
		appendCondition(fmt.Sprintf("match_id = $%d", counter), "AND", options.Match)
	}

	if options.Team != "" {
		appendCondition(
			fmt.Sprintf("(home_team_id = $%d OR away_team_id = $%d)", counter, counter),
			"AND",
			options.Team,
		)
	}

	if options.Username != "" {
		appendCondition(
			fmt.Sprintf("username = $%d", counter),
			"AND",
			options.Username,
		)
	}

	if options.SearchTerm != "" {
		appendCondition(
			fmt.Sprintf("r.search_vector @@ websearch_to_tsquery('portuguese', $%d)", counter),
			"AND",
			options.SearchTerm,
		)
	}

	if isFirstCondition {
		return "", make([]any, 0)
	}

	return "WHERE " + whereClause.String(), params
}
