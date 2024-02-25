package pkmn_errors

import "net/http"

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int, msg string) {
	switch status {
	case http.StatusNotFound:
		notFoundError(w, r, msg)
	case http.StatusInternalServerError:
		internalServerError(w, r, msg)
	}
}

func internalServerError(w http.ResponseWriter, r *http.Request, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error: " + msg))
}

func notFoundError(w http.ResponseWriter, r *http.Request, msg string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found: " + msg))
}
