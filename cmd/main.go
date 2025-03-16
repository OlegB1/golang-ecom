package main

import (
	"log"

	"github.com/OlegB1/ecom/cmd/api"
	"github.com/OlegB1/ecom/config"
	"github.com/OlegB1/ecom/db"
)

func main() {
	db := db.NewStorage(config.Envs.DB_ADDR)
	port := config.Envs.SERVER_ADDR
	if port == "" {
		port = "8080"
	}

	server := api.NewAPIServer(":"+port, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
