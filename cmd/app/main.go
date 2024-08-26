package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jalenarms1/sillysocks-GoTH/db"
	"github.com/Jalenarms1/sillysocks-GoTH/handlers"
	"github.com/Jalenarms1/sillysocks-GoTH/models"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	gob.Register(&models.Cart{})
	gob.Register(&models.CartItem{})
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hello World")

	db.InitDB(os.Getenv("DB_CONN_STR"))
	defer db.DB.Close()

	router := chi.NewMux()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	FileServer(router, "/public", filesDir)

	listenAddr := os.Getenv("LISTEN_ADDR")

	handlers.RegisterRouter(router)

	router.Get("/", handlers.UseHTTPHandler(handlers.HandleRoot))
	fmt.Printf("http://localhost%s\n", listenAddr)
	http.ListenAndServe(listenAddr, router)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
