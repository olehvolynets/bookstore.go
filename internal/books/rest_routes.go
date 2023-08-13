package books

import (
	"net/http"

	"bookstore/internal/logging"
)

func HandleGetBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	w.WriteHeader(http.StatusAccepted)
	_, err := w.Write([]byte("{\"data\":[]}"))
	if err != nil {
		log.Error().Err(err).Msg("error while writing a response")
	}
}

func HandleGetBook(w http.ResponseWriter, _ *http.Request) {
}
