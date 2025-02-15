package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jalenarms1/sillysocks-GoTH/internal/db"
	"github.com/stripe/stripe-go/v81"
)

func init() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("ENV not loaded")
	// }

	if err := db.SetDB(); err != nil {
		log.Fatal(err)
	}

	stripe.Key = os.Getenv("STRIPE_SKEY")

	fmt.Println("Connected to DB")

}

func main() {
	mux := http.NewServeMux()

	handler := registerRoutes(mux)

	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDR"), handler))
}
