package api

import (
	db "SimpleBank/db/sqlc"
	"SimpleBank/db/util"
	"SimpleBank/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", ValidCurrency)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("business", ValidBusiness)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	{
		authRoutes.POST("/accounts", server.createAccount)
		authRoutes.GET("/accounts/:id", server.getAccount)
		authRoutes.GET("/accounts", server.listAccount)
		authRoutes.POST("/accounts/update", server.updateAccount)
		authRoutes.DELETE("/accounts/:id", server.deleteAccount)

		authRoutes.POST("/transfers", server.createTransfer)

		authRoutes.POST("/business", server.NewBusiness)
		authRoutes.GET("/business", server.listEntries)
	}

	server.router = router
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
