package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Jalenarms1/sillysocks-GoTH/internal/db"
)

func HandleGetOrder(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	orderId := r.URL.Path[len("/order"):]

	order := db.GetOrder((orderId))
	if order == nil {
		return errors.New("order not found " + orderId)

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)

	return nil
}
