package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

func RegisterAuth(r *chi.Mux) {
	r.Get("/api/user/token", UseHTTPHandler(initUserToken))
}

func initUserToken(w http.ResponseWriter, r *http.Request) error {

	token, _ := uuid.NewV4()

	data := []byte(fmt.Sprintf(`{
		"token": "%s"
	}`, token))

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}
