package errors

type ErrNotFound struct {
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
