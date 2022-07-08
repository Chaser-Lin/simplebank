package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testing"
)

const (
	dbDriver = "mysql"
	dbSource = "root:secret@tcp(localhost:3306)/simple_bank?parseTime=true"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
