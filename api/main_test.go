package api

import (
	db "SimpleBank/db/sqlc"
	util2 "SimpleBank/util"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const (
	dbDriver = "mysql"
	dbSource = "root:secret@tcp(localhost:3306)/simple_bank?parseTime=true"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util2.Config{
		TokenSymmetricKey:   util2.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
