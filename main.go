package main

import (
	"SimpleBank/api"
	db "SimpleBank/db/sqlc"
	"SimpleBank/db/util"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//const (
//	dbDriver      = "mysql"
//	dbSource      = "root:secret@tcp(localhost:3306)/simple_bank?parseTime=true"
//	serverAddress = "0.0.0.0:8080"
//)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load congfig:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
