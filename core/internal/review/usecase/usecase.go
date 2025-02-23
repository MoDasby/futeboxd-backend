package usecase

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/modasby/futeboxd-backend/core/internal/review"
	"github.com/modasby/futeboxd-backend/core/internal/review/dto"
	"github.com/modasby/futeboxd-backend/core/internal/user"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/football"
	"github.com/modasby/futeboxd-backend/core/pkg/json/null"
	"github.com/modasby/futeboxd-backend/core/pkg/pagination"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type reviewUsecases struct {
	reviewRepo     review.Repository
	footballClient football.Client
	userRepo       user.Repository
}

func NewReviewUsecases(
	reviewRepo review.Repository,
	footballClient football.Client,
	userRepo user.Repository,
) review.ReviewUsecases {
	return &reviewUsecases{
		reviewRepo:     reviewRepo,
		footballClient: footballClient,
		userRepo:       userRepo,
	}
}

func (uc *reviewUsecases) Create(ctx context.Context, input *dto.ReviewInput) error {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	match, err := uc.footballClient.GetMatch(ctx, input.MatchID)
	if err != nil {
		return err
	}

	author, err := uc.userRepo.FindOneByIdOrUsername(ctx, session.UserID)
	if err != nil {
		return err
	}

	review, err := review.NewReview(
		author,
		input.Rate,
		match.ID,
		input.Description,
		match.HomeCompetitor.Team.ID,
		match.AwayCompetitor.Team.ID,
	)
	if err != nil {
		return err
	}

	err = uc.reviewRepo.Upsert(ctx, review)
	if err != nil {
		return err
	}

	return nil
}

func (uc *reviewUsecases) Delete(ctx context.Context, reviewID int64) error {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return err
	}

	review, err := uc.reviewRepo.FindOneByID(ctx, session.UserID, reviewID)
	if err != nil {
		return err
	}

	if review.Author.ID != session.UserID {
		return &errors.HTTPErr{
			Msg:        "você não pode executar essa ação",
			Code:       http.StatusForbidden,
			Context:    "REVIEW:USECASE:DELETE_REVIEW:FORBIDDEN",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}
	}

	if err := uc.reviewRepo.Delete(ctx, reviewID); err != nil {
		return err
	}

	return nil
}

func (uc *reviewUsecases) ListFeed(ctx context.Context, page *pagination.Page) ([]dto.Review, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	reviews, err := uc.reviewRepo.ListFeed(ctx, session.UserID, page)
	if err != nil {
		return nil, err
	}

	matches, err := uc.getMatchesByReview(ctx, uc.footballClient, reviews)
	if err != nil {
		return nil, err
	}

	output := make([]dto.Review, 0)

	for _, review := range reviews {
		match := matches[review.MatchID]

		output = append(output, dto.Review{
			ID: review.ID,
			Author: dto.ContentAuthor{
				ID:             review.Author.ID,
				Name:           null.String(review.Author.Name),
				Username:       review.Author.Username,
				ProfilePicture: review.Author.ProfilePicture,
			},
			Rate:          review.Rate,
			Description:   review.Description,
			Match:         &match,
			Likes:         review.Likes,
			CommentsCount: review.CommentsCount,
			IsLiked:       review.IsLiked,
			CreatedAt:     review.CreatedAt,
		})
	}

	return output, nil
}

func (uc *reviewUsecases) ListBy(
	ctx context.Context,
	options review.ListReviewsOptions,
	page *pagination.Page,
) ([]dto.Review, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	whereClause, params := uc.buildWhereClause(options)

	reviews, err := uc.reviewRepo.ListAll(ctx, session.UserID, whereClause, params, page)
	if err != nil {
		return nil, err
	}

	matches, err := uc.getMatchesByReview(ctx, uc.footballClient, reviews)
	if err != nil {
		return nil, err
	}

	output := make([]dto.Review, 0)

	for _, review := range reviews {
		match := matches[review.MatchID]

		output = append(output, dto.Review{
			ID: review.ID,
			Author: dto.ContentAuthor{
				ID:             review.Author.ID,
				Name:           null.String(review.Author.Name),
				Username:       review.Author.Username,
				ProfilePicture: review.Author.ProfilePicture,
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

func (uc *reviewUsecases) Search(ctx context.Context, term string, page *pagination.Page) ([]dto.Review, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	reviews, err := uc.reviewRepo.Search(ctx, session.UserID, term, page)
	if err != nil {
		return nil, err
	}

	matches := make(map[int64]football.Match, 0)

	if len(reviews) > 0 {
		matches, err = uc.getMatchesByReview(ctx, uc.footballClient, reviews)
		if err != nil {
			return nil, err
		}
	}

	output := make([]dto.Review, len(reviews))

	for index, review := range reviews {
		match := matches[review.MatchID]

		output[index] = dto.Review{
			ID: review.ID,
			Author: dto.ContentAuthor{
				ID:             review.Author.ID,
				Name:           null.String(review.Author.Name),
				Username:       review.Author.Username,
				ProfilePicture: review.Author.ProfilePicture,
			},
			Rate:          review.Rate,
			Description:   review.Description,
			Match:         &match,
			CommentsCount: review.CommentsCount,
			Likes:         review.Likes,
			IsLiked:       review.IsLiked,
			CreatedAt:     review.CreatedAt,
		}
	}

	return output, nil
}

func (uc *reviewUsecases) ToggleLike(ctx context.Context, reviewID int64) (*dto.LikeStats, error) {
	session, err := utils.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	reviewExists, err := uc.reviewRepo.ExistsByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}

	if !reviewExists {
		httpErr := &errors.HTTPErr{
			Msg:        "essa review não existe",
			Code:       http.StatusNotFound,
			Context:    "REVIEW:USECASE:TOGGLE_LIKE_REVIEW:REVIEW_NOT_FOUND",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}

		return nil, httpErr
	}

	isLiked, err := uc.reviewRepo.IsLiked(ctx, session.UserID, reviewID)
	if err != nil {
		return nil, err
	}

	if isLiked {
		if err := uc.reviewRepo.Unlike(ctx, session.UserID, reviewID); err != nil {
			return nil, err
		}

		return &dto.LikeStats{Like: false}, nil
	}

	if err := uc.reviewRepo.Like(ctx, session.UserID, reviewID); err != nil {
		return nil, err
	}

	return &dto.LikeStats{Like: true}, nil
}

func (uc *reviewUsecases) getMatchesByReview(ctx context.Context, footballClient football.Client, reviews []review.Review) (map[int64]football.Match, error) {
	matchIDs := make([]int64, len(reviews))

	for i, review := range reviews {
		matchIDs[i] = review.MatchID
	}

	return footballClient.GetMatchesMap(ctx, matchIDs)
}

func (uc *reviewUsecases) buildWhereClause(options review.ListReviewsOptions) (string, []any) {
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
