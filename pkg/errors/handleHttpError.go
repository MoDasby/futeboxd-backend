package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

type responseBody struct {
	Msg  string `json:"message"`
	Code int32  `json:"code"`
}

func HandleHttpError(w http.ResponseWriter, err error) {

	log.Printf("ERROR: %s", err.Error())

	body := responseBody{Msg: err.Error()}

	switch err.(type) {
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
	default:
		w.WriteHeader(500)
	}
}
