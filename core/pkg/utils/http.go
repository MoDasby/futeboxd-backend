package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/modasby/futeboxd-backend/core/internal/session"
	"github.com/modasby/futeboxd-backend/core/pkg/errors"
	"github.com/modasby/futeboxd-backend/core/pkg/middleware"
)

func ParseIntValue(ctx context.Context, value string) (int64, error) {
	valueInt, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, &errors.HTTPErr{
			Msg:        "Valor inválido ",
			Code:       http.StatusBadRequest,
			Context:    "UTILS:PARSE_INT_VALUE:", // incompleto para o handler completar
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}
	}

	return int64(valueInt), nil
}

func GetTraceIDFromCtx(ctx context.Context) string {
	traceID := ctx.Value("traceID")
	if traceID == nil {
		return ""
	}

	return traceID.(string)
}

func GetSessionFromCtx(ctx context.Context) (*session.Session, error) {
	session, ok := ctx.Value(middleware.SessionKey).(*session.Session)
	if !ok || session == nil {
		return nil, &errors.HTTPErr{
			Msg:        "Sessão inválida",
			Code:       http.StatusBadRequest,
			Context:    "UTILS:GET_SESSION_FROM_CTX:INVALID_SESSION",
			StackTrace: errors.CaptureStackTrace(),
			ErrorCode:  GetTraceIDFromCtx(ctx),
			Timestamp:  time.Now().UTC(),
		}
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
