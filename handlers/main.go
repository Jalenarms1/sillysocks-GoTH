package handlers

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRouter(router *chi.Mux) {
	RegisterProductRouter(router)
	// RegisterCartRouter(router)
	RegisterStripeRouter(router)
	RegisterAuth(router)
	RegisterAdmin(router)
}
