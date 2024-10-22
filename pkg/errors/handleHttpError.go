package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

type responseBody struct {
	Msg  string `json:"message"`
	Code int    `json:"code"`
}

func HandleHttpError(w http.ResponseWriter, err error) {

	log.Println(err.Error())

	if e, ok := err.(*HTTPErr); ok {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(e.Code)

		body := responseBody{
			Msg:  e.Msg,
			Code: e.Code,
		}

		json.NewEncoder(w).Encode(body)

		return
	}

	w.WriteHeader(500)

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
