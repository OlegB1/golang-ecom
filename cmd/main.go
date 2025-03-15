package main

import (
	"log"

	"github.com/OlegB1/ecom/cmd/api"
	"github.com/OlegB1/ecom/config"
	"github.com/OlegB1/ecom/db"
)

func main() {
	db := db.NewStorage(config.Envs.DB_ADDR)

	server := api.NewAPIServer(":"+config.Envs.SERVER_ADDR, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
