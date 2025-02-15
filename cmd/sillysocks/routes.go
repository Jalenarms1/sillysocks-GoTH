package main

import (
	"net/http"

	"github.com/Jalenarms1/sillysocks-GoTH/internal/handlers"
)

func registerRoutes(mux *http.ServeMux) http.Handler {
	mux.HandleFunc("/products", handlers.HandleGetProducts)
	mux.HandleFunc("/products/{id}", handlers.HandleGetProduct)

	mux.HandleFunc("/checkout", handlers.HandleCreateCheckoutSession)
	mux.HandleFunc("/checkout-wh", handlers.HandleCheckoutSessionWebhook)

	wrappedMux := handlers.UseCors(mux)

	return wrappedMux
}
