package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Jalenarms1/sillysocks-GoTH/internal/db"
)

func HandleGetProducts(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}
	products := db.GetProducts()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(products)

	return nil

}

func HandleGetProduct(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	id := r.URL.Path[len("/products/"):]
	product := db.GetProduct(id)

	if product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return errors.New("product not found")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(product)
	return nil
}
