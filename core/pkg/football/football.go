package football

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/modasby/futeboxd-backend/core/config"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/utils"
)

type footballClient struct {
	baseUrl string
}

func NewClient(cfg config.Football) Client {
	return &footballClient{
		baseUrl: cfg.Url,
	}
}

func (fs *footballClient) GetMatchesMap(ctx context.Context, matchIDs []int64) (map[int64]Match, error) {
	matches, err := fs.GetMatches(ctx, matchIDs)
	if err != nil {
		return nil, err
	}

	matchesMap := make(map[int64]Match)

	for _, match := range matches {
		matchesMap[match.ID] = match
	}

	return matchesMap, nil
}

func (fs *footballClient) GetMatches(ctx context.Context, matchIDs []int64) ([]Match, error) {
	body, err := json.Marshal(matchIDs)
	if err != nil {
		return nil, err
	}

	// TODO setar o header User-Agent pra core client

	res, err := http.Post(fmt.Sprintf("%s/matches/batch", fs.baseUrl), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var body errors.ResponseBody

		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}

		return nil, &errors.HTTPErr{
			Msg:        body.Msg,
			Code:       res.StatusCode,
			Context:    "FOOTBALL:CLIENT:GET_MATCHES",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
			Original:   err,
		}
	}

	matches := make([]Match, 0)

	if err := json.NewDecoder(res.Body).Decode(&matches); err != nil {
		return nil, err
	}

	return matches, err
}

func (fs *footballClient) GetMatch(ctx context.Context, matchID int64) (*Match, error) {
	res, err := http.Get(fmt.Sprintf("%s/matches/%d", fs.baseUrl, matchID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var body errors.ResponseBody

		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}

		return nil, &errors.HTTPErr{
			Msg:        body.Msg,
			Code:       res.StatusCode,
			Context:    "FOOTBALL:CLIENT:GET_MATCH",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
			Original:   err,
		}
	}
	var match Match

	if err := json.NewDecoder(res.Body).Decode(&match); err != nil {
		return nil, err
	}

	return &match, err
}

func (fs *footballClient) GetTeam(ctx context.Context, teamID int64) (*Team, error) {
	res, err := http.Get(fmt.Sprintf("%s/teams/%d", fs.baseUrl, teamID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var body errors.ResponseBody

		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}

		return nil, &errors.HTTPErr{
			Msg:        body.Msg,
			Code:       res.StatusCode,
			Context:    "FOOTBALL:CLIENT:GET_TEAM",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  utils.GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
			Original:   err,
		}
	}
	var team Team

	if err := json.NewDecoder(res.Body).Decode(&team); err != nil {
		return nil, err
	}

	return &team, err
}
