package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HttpError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ResponseData struct {
	Message   string      `json:"message"`
	Status    int         `json:"status"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
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

func ResponseOk(w http.ResponseWriter, message string, payload interface{}) {
	w.WriteHeader(200)
	out := ResponseData{Message: message, Status: 200}
	out.Payload = payload

	json.
		NewEncoder(w).
		Encode(out)
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
