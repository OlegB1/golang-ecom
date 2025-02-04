package main

import (
	"database/sql"
	"log"

	"github.com/OlegB1/ecom/cmd/api"
	"github.com/OlegB1/ecom/config"
	"github.com/OlegB1/ecom/db"
)

func main() {
	db, err := db.NewStorage(config.Envs.DB_ARRD)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":"+config.Envs.SERVER_ARRD, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected")
}
