package handlers

import (
	"net/http"

	"github.com/Jalenarms1/sillysocks-GoTH/models"
	"github.com/Jalenarms1/sillysocks-GoTH/views/home"
	"github.com/go-chi/chi/v5"
)

func handleGetProducts(w http.ResponseWriter, r *http.Request) error {
	ps, err := models.GetProducts()
	if err != nil {
		return err
	}
	// cart := models.GetSessionValues(r)
	// cart.Values["cart"] = models.Cart{}
	// saveErr := cart.Save(r, w)
	// if saveErr != nil {
	// 	log.Fatal((saveErr))
	// }
	// fmt.Print(cart.Values["cart"])
	Render(w, r, home.CatalogSlides(ps))

	return nil
}

func RegisterProductRouter(router *chi.Mux) {
	router.Get("/api/products/list", UseHTTPHandler(handleGetProducts))
}
