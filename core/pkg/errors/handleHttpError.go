package errors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func HandleHttpError(w http.ResponseWriter, err error) {

	if e, ok := err.(*HTTPErr); ok {
		sendResponse(w, e)

		log.Printf("CORE: %s", err.Error())

		return
	}

	httpErr := NewHTTPErr(
		"ocorreu um erro desconhecido",
		500,
		fmt.Sprintf("CORE: INTERNAL_SERVER_ERROR: %s", err.Error()),
	)

	sendResponse(w, httpErr)

	log.Println(httpErr.Error())

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
