package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewStorage(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected")
	return db
}
