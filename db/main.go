package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(connString string) {
	var err error
	DB, err = sql.Open("postgres", connString)

	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
