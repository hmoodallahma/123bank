package main

import (
	"database/sql"
	"github.com/hmoodallahma/123bank/api"
	db "github.com/hmoodallahma/123bank/db/sqlc"
	"github.com/hmoodallahma/123bank/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configurations", err)
	}

	// connect to db and create a store
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
