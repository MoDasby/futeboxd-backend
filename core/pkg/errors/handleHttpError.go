package errors

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime"
	"time"
)

type ResponseBody struct {
	Msg  string `json:"message"`
	Code string `json:"code"`
}

func sendResponse(w http.ResponseWriter, err *HTTPErr) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	body := ResponseBody{
		Msg:  err.Msg,
		Code: err.ErrorCode,
	}

	json.NewEncoder(w).Encode(body)
}

func CaptureStackTrace() string {
	stackBuf := make([]byte, 1024)
	stackLen := runtime.Stack(stackBuf, false)
	return string(stackBuf[:stackLen])
}

func LogError(ctx context.Context, err *HTTPErr) {
	originalErrMsg := ""

	if err.Original != nil {
		originalErrMsg = err.Original.Error()
	}

	slog.ErrorContext(
		ctx,
		err.Msg,
		"code", err.Code,
		"stackTrace", err.StackTrace,
		"context", err.Context,
		"traceID", err.ErrorCode,
		"timestamp", err.Timestamp,
		"originalError", originalErrMsg,
	)
}

func HandleHttpError(ctx context.Context, w http.ResponseWriter, err error) {
	if e, ok := err.(*HTTPErr); ok {
		sendResponse(w, e)

		LogError(ctx, e)

		return
	}

	traceID, ok := ctx.Value("traceID").(string)
	if !ok {
		traceID = ""
	}

	httpErr := &HTTPErr{
		Msg:        "Ocorreu um erro desconhecido",
		Code:       http.StatusInternalServerError,
		StackTrace: CaptureStackTrace(),
		Context:    "UNKNOWN_ERROR",
		ErrorCode:  traceID,
		Timestamp:  time.Now().UTC(),
		Original:   err,
	}

	sendResponse(w, httpErr)

	LogError(ctx, httpErr)

	/* switch err.(type) {
	case *ErrNotFound:
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		body.Code = http.StatusNotFound

		json.NewEncoder(w).Encode(body)
	case *ErrBadRequest:
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		body.Code = http.StatusBadRequest

		json.NewEncoder(w).Encode(body)
	case *ErrUnauthorized:
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		body.Code = http.StatusUnauthorized

		json.NewEncoder(w).Encode(body)
	case *ErrForbidden:
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		body.Code = http.StatusForbidden

		json.NewEncoder(w).Encode(body)
	case *ErrInternalServer:
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		body.Code = http.StatusInternalServerError

		json.NewEncoder(w).Encode(body)
	default:
		w.WriteHeader(500)
	} */
}
