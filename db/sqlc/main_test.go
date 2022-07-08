package db

import (
	"SimpleBank/db/util"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "mysql"
	dbSource = "root:secret@tcp(localhost:3306)/simple_bank?parseTime=true"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	//conn, err := sql.Open(dbDriver, dbSource)
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load congfig: ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	//testQueries = New(conn)
	testQueries = New(testDB)
	os.Exit(m.Run())
}
