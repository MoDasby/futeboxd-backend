package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/modasby/futeboxd-api/services/core/internal/domain"
	"github.com/modasby/futeboxd-api/services/core/internal/middleware"
)

func ParseIntValue(value string) (int64, error) {
	valueInt, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, errors.New("invalid int value")
	}

	return int64(valueInt), nil
}

func GetUserFromCtx(ctx context.Context) (*domain.User, error) {
	user, ok := ctx.Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil {
		return nil, errors.New("can't parse user")
	}

	return user, nil
}

func sendJsonResponse(w http.ResponseWriter, res any) error {
	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		return err
	}

	return nil
}
