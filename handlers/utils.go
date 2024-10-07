package handlers

import (
	"log/slog"
	"net/http"
	"os"
)

type UserContextKey string

const (
	UserCtxKey UserContextKey = "user-uid-silly-socks"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

func UseHTTPHandler(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("HTTP Error", "err", err, "path", r.URL)
		}
	}
}

func UseAdminHTTPHandler(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(UserCtxKey) != os.Getenv("SILLYSOCKS_ADMIN_KEY") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := h(w, r); err != nil {
			slog.Error("HTTP Error", "err", err, "path", r.URL)
		}
	}
}
