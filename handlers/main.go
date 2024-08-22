package handlers

import (
	"net/http"

	"github.com/Jalenarms1/sillysocks-GoTH/views"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, views.Home("Jalen"))
}
