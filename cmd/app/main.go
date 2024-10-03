package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jalenarms1/sillysocks-GoTH/db"
	"github.com/Jalenarms1/sillysocks-GoTH/handlers"
	"github.com/go-chi/chi/v5"
)

func userMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CLIENT_DOMAIN"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// mh := r.Header["X-silly-socks-app"]

		// if mh != nil {
		// 	isMobile := mh[0] == os.Getenv("MOBILE_TOKEN")
		// }

		token := r.Header.Get("Authorization")
		var ctx context.Context
		if token != "" || r.URL.Path == "/api/user/token" {
			fmt.Println("Token found", token)
			// http.SetCookie(w, &http.Cookie{
			// 	Name:     "silly-socks-user",
			// 	Value:    localId.String(),
			// 	Path:     "/",
			// 	Secure:   false,
			// 	HttpOnly: true,
			// 	SameSite: http.SameSiteNoneMode,
			// })

			ctx = context.WithValue(r.Context(), handlers.UserCtxKey, token)

		} else {
			fmt.Println("not found")
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {

	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal(err)
	// }

	db.InitDB(os.Getenv("MASTER_DB_URL"))
	defer db.CloseDB()

	router := chi.NewMux()

	router.Use(userMiddleware)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	FileServer(router, "/public", filesDir)
	listenAddr := os.Getenv("LISTEN_ADDR")
	fmt.Printf("%s\n", listenAddr)

	handlers.RegisterRouter(router)

	// router.Get("/", handlers.UseHTTPHandler(handlers.HandleRoot))
	fmt.Printf("http://localhost%s\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, router))
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
