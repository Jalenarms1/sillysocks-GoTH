package handlers

import (
	"fmt"
	"net/http"
)

func UseCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})

}

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func ErrorCatchHandler(fn ErrorHandlerFunc) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)

		}

	}

}
