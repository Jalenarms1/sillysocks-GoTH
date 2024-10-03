package handlers

import (
	"log/slog"
	"net/http"
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
