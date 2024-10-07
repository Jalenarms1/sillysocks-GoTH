package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Jalenarms1/sillysocks-GoTH/db"
	"github.com/Jalenarms1/sillysocks-GoTH/models"
	"github.com/go-chi/chi/v5"
)

func RegisterAdmin(r *chi.Mux) {
	r.Delete("/api/admin/delete-products", UseAdminHTTPHandler(deleteProducts))

	r.Post("/api/admin/add-product", UseAdminHTTPHandler(addProduct))
}

func addProduct(w http.ResponseWriter, r *http.Request) error {

	// body, _ := io.ReadAll(r.Body)

	var product models.ProductFormReq

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		fmt.Println("err here")
		return err
	}

	if iErr := models.SubmitProduct(&product); iErr != nil {
		return iErr
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func deleteProducts(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var data []string

	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		return jsonErr
	}

	fmt.Println(data)

	var formattedQStrs []string
	for _, s := range data {
		formattedQStrs = append(formattedQStrs, fmt.Sprintf(`'%s'`, s))
	}

	query := fmt.Sprintf(`
		DELETE FROM "Product" where "Id" in (%s)
	`, strings.Join(formattedQStrs, ","))

	_, exErr := db.DB.Exec(query)
	if exErr != nil {
		return exErr
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
