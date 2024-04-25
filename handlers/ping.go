package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func PingHandler(r *mux.Router) {
	r.HandleFunc(
		"/ping",
		func(w http.ResponseWriter, r *http.Request) {
			SendResultResponse(w, http.StatusOK, nil)
			return
		},
	).Methods(http.MethodGet)
}
