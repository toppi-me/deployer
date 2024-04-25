package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/toppi-me/deployer/internal/log"
)

type Response struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *ResponseError  `json:"error,omitempty"`
}

type ResponseError struct {
	Message string `json:"message"`
}

func (r Response) Send(w http.ResponseWriter, HTTPCode int) {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(HTTPCode)

	responseJSON, err := json.Marshal(r)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	_, err = w.Write(responseJSON)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
}

func SendErrorResponse(w http.ResponseWriter, httpCode int, message string) {
	var response Response

	response.Ok = false
	response.Error = &ResponseError{
		Message: message,
	}

	response.Send(w, httpCode)
}

func SendResultResponse(w http.ResponseWriter, httpCode int, result json.RawMessage) {
	var response Response

	response.Ok = true
	response.Result = result

	response.Send(w, httpCode)
}
