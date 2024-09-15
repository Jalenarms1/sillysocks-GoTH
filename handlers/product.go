package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Jalenarms1/sillysocks-GoTH/db"
	"github.com/Jalenarms1/sillysocks-GoTH/models"
	"github.com/go-chi/chi/v5"
)

func RegisterProductRouter(r *chi.Mux) {
	r.Get("/api/products/list", UseHTTPHandler(handleGetProductList))
	r.Get("/api/products/find/{productId}", UseHTTPHandler(handleGetProductById))
}

func handleGetProductById(w http.ResponseWriter, r *http.Request) error {
	productId := chi.URLParam(r, "productId")

	var product models.Product

	if err := db.DB.Get(&product, `SELECT * FROM Product WHERE Id = ?`, productId); err != nil {
		return err
	}

	jErr := json.NewEncoder(w).Encode(product)
	if jErr != nil {
		return jErr
	}

	return nil
}

func handleGetProductList(w http.ResponseWriter, r *http.Request) error {
	var products []models.Product

	// if err := db.DB.Select(&products, `SELECT * FROM Product`); err != nil {
	// 	return err
	// }

	url := "https://masterdb-jalenarms1.turso.io/v2/pipeline"

	jErr := json.NewEncoder(w).Encode(products)
	if jErr != nil {
		return jErr
	}

	return nil
}
