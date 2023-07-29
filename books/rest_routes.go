package books

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func HandleGetBooks(w http.ResponseWriter, _ bunrouter.Request) error {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("{\"data\":[]}"))

	return nil
}
