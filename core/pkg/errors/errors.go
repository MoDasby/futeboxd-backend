package errors

import (
	"fmt"
	"time"
)

type HTTPErr struct {
	Msg        string
	Code       int
	StackTrace string
	Context    string
	ErrorCode  string
	Timestamp  time.Time
	Original   error
}

func (err HTTPErr) Error() string {
	var originalErr error

	if err.Original != nil {
		originalErr = err.Original
	}

	return fmt.Sprintf("Error: %s, HttpCode: %d, Context: %s, ErrorCode: %s, StackTrace: %s, Timestamp: %s, Original: %s\n",
		err.Msg, err.Code, err.Context, err.ErrorCode,
		err.StackTrace, err.Timestamp.String(), originalErr,
	)
}

/* type ErrNotFound struct {
	Msg string
}

func (err *ErrNotFound) Error() string {
	return err.Msg
}

func NewErrNotFound(msg string) *ErrNotFound {
	return &ErrNotFound{Msg: msg}
}

// -----

type ErrBadRequest struct {
	Msg string
}

func (err *ErrBadRequest) Error() string {
	return err.Msg
}

func NewErrBadRequest(msg string) *ErrBadRequest {
	return &ErrBadRequest{Msg: msg}
}

// -----

type ErrUnauthorized struct {
	Msg string
}

func (err *ErrUnauthorized) Error() string {
	return err.Msg
}

func NewErrUnauthorized(msg string) *ErrUnauthorized {
	return &ErrUnauthorized{Msg: msg}
}

// ------

type ErrForbidden struct {
	Msg string
}

func (err *ErrForbidden) Error() string {
	return err.Msg
}

func NewErrForbidden(msg string) *ErrForbidden {
	return &ErrForbidden{Msg: msg}
}

// -----

type ErrInternalServer struct {
	Msg string
}

func (err *ErrInternalServer) Error() string {
	return err.Msg
}

func NewErrInternalServer(msg string) *ErrInternalServer {
	return &ErrInternalServer{Msg: msg}
}
*/
