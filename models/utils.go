package models

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

func generateUUIDv4() (string, error) {
	var uuid [16]byte
	_, err := rand.Read(uuid[:])
	if err != nil {
		return "", err
	}

	// Set version (4 bits) to 0100
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set variant (2 bits) to 10
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	// Format UUID as a string
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			uuid[0:4],
			uuid[4:6],
			uuid[6:8],
			uuid[8:10],
			uuid[10:]),
		nil
}

func DataGateway[T any](f func() ([]T, error)) []T {
	data, err := f()
	if err != nil {
		log.Fatal(err)
	}

	return data
}

// var store = sessions.NewCookieStore([]byte(os.Getenv("SILLYSOCKS_SESSION_KEY")))

func GetSessionValues(r *http.Request) *sessions.Session {
	var store = sessions.NewCookieStore([]byte(os.Getenv("SILLYSOCKS_SESSION_KEY")))

	fmt.Println(os.Getenv("SILLYSOCKS_SESSION_KEY"))
	session, err := store.Get(r, "sillysocks_cart")
	if err != nil {
		log.Fatal(err)
	}

	return session
}
