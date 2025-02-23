package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func LoggerMiddleware(ctx context.Context, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()

		ctx = context.WithValue(ctx, "traceID", traceID)

		slog.InfoContext(r.Context(), "Request recebida",
			"method", r.Method,
			"path", r.URL.Path,
			"traceID", traceID,
			"clientIP", r.Header.Get("X-Forwarded-For"),
		)

		startTime := time.Now()

		next(w, r.WithContext(ctx))

		duration := time.Since(startTime)

		slog.Info("Request processada",
			"method", r.Method,
			"path", r.URL.Path,
			"traceID", traceID,
			"clientIP", r.Header.Get("X-Forwarded-For"),
			"durationMs", duration.Milliseconds(),
		)
	}
}
