package books

import (
	"net/http"
)

func HandleGetBooks(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("{\"data\":[]}"))
}

func HandleGetBook(w http.ResponseWriter, _ *http.Request) {
}
