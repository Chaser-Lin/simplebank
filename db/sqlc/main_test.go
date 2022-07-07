package db

import (
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
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("db connect error: ", err)
	}
	//testQueries = New(conn)
	testQueries = New(testDB)
	os.Exit(m.Run())
}
