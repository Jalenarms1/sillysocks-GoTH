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

	orderId := r.URL.Path[len("/order/"):]

	order, err := db.GetOrder((orderId))
	if err != nil {
		return errors.New("order not found " + orderId + " " + err.Error())

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)

	return nil
}
