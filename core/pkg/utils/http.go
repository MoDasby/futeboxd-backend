package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
)

func ParseIntValue(value string) (int64, error) {
	valueInt, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, errors.NewHTTPErr(
			"valor inválido ",
			400,
			"UTILS:PARSE_INT_VALUE:", // incompleto para o handler completar
		)
	}

	return int64(valueInt), nil
}

func GetSessionFromCtx(ctx context.Context) (*session.Session, error) {
	session, ok := ctx.Value(middleware.SessionKey).(*session.Session)
	if !ok || session == nil {
		return nil, errors.NewHTTPErr(
			"sessão inválida",
			401,
			"UTILS:GET_SESSION_FROM_CTX:INVALID_SESSION",
		)
	}

	return session, nil
}

func SendJSON(w http.ResponseWriter, res any) error {
	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		return err
	}

	return nil
}
