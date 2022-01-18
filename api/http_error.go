package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func newHttpError(w http.ResponseWriter, code int, msg string, err error) {
	SetResponseHeader(w)
	w.WriteHeader(code)
	body := &HttpError{
		Code:    code,
		Message: msg,
		Error:   fmt.Sprintf("%+v", err),
	}
	json.
		NewEncoder(w).
		Encode(body)

}

func ResponseOk(w http.ResponseWriter) {
	w.WriteHeader(200)
}

func ResponseError(w http.ResponseWriter, code int, msg string, err error) {
	newHttpError(w, code, msg, err)
}

func ResponseErrorForbidden(w http.ResponseWriter, err error) {
	code := http.StatusForbidden
	msg := http.StatusText(http.StatusForbidden)
	newHttpError(w, code, msg, err)

}

/*
 * response code 401 StatusUnauthorized
 */
func ResponseErrorUnauthorized(w http.ResponseWriter, err error) {
	code := http.StatusUnauthorized
	msg := http.StatusText(http.StatusUnauthorized)
	newHttpError(w, code, msg, err)
}
