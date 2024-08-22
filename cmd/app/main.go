package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jalenarms1/sillysocks-GoTH/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hello World")

	router := chi.NewMux()

	listenAddr := os.Getenv("LISTEN_ADDR")

	router.Get("/", handlers.Make(handlers.HandleRoot))
	fmt.Printf("http://localhost%s\n", listenAddr)
	http.ListenAndServe(listenAddr, router)
}
