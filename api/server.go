package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/hmoodallahma/123bank/db/sqlc"
	"github.com/hmoodallahma/123bank/token"
	"github.com/hmoodallahma/123bank/util"
)

// Servers HTTP requests for banking service
type Server struct {
	config     util.Config
	store      db.Store    // access db
	tokenMaker token.Maker // make auth tokens
	router     *gin.Engine // handles routing
}

// creates a new server instance
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// add validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	// add routes to router
	server.setUpRouter()

	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()
	router.POST("/users/login", server.loginUser)
	router.POST("/user", server.createUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/account/:id", server.getAccount)
	authRoutes.PUT("/account/:id", server.updateAccount)
	authRoutes.POST("/account", server.createAccount)
	authRoutes.DELETE("/account/:id", server.deleteAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.DELETE("/user/:username", server.deleteUser)
	authRoutes.GET("/users", server.listUsers)

	authRoutes.POST("/transfer", server.createTransfer)
	server.router = router
}

// start the server on a specific address
func (server *Server) Start(address string) error {
	// router is private and cannot be accessed outside of pacakge
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
