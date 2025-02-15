package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var DB *sql.DB

func SetDB() error {
	var err error
	DB, err = sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", os.Getenv("DB_URL"), os.Getenv("DB_TOKEN")))

	return err
}
