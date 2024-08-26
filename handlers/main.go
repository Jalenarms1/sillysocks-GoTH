package handlers

import (
	"net/http"

	"github.com/Jalenarms1/sillysocks-GoTH/views/home"
	"github.com/go-chi/chi/v5"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, home.Page("Jalen"))
}

func RegisterRouter(router *chi.Mux) {
	RegisterProductRouter(router)
	RegisterCartRouter(router)
}
