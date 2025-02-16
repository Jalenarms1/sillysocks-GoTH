package main

import (
	"net/http"

	"github.com/Jalenarms1/sillysocks-GoTH/internal/handlers"
)

func registerRoutes(mux *http.ServeMux) http.Handler {
	mux.HandleFunc("/products", handlers.ErrorCatchHandler(handlers.HandleGetProducts))
	mux.HandleFunc("/products/{id}", handlers.ErrorCatchHandler(handlers.HandleGetProduct))

	mux.HandleFunc("/checkout", handlers.ErrorCatchHandler(handlers.HandleCreateCheckoutSession))
	mux.HandleFunc("/checkout-wh", handlers.ErrorCatchHandler(handlers.HandleCheckoutSessionWebhook))

	mux.HandleFunc("/order/{id}", handlers.ErrorCatchHandler(handlers.HandleGetOrder))

	wrappedMux := handlers.UseCors(mux)

	return wrappedMux
}
