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
	router := gin.Default()

	// add validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	// add routes to router
	router.POST("/account", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.PUT("/account/:id", server.updateAccount)
	router.DELETE("/account/:id", server.deleteAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfer", server.createTransfer)
	router.POST("/user", server.createUser)
	router.DELETE("/user/:username", server.deleteUser)
	router.GET("/users", server.listUsers)
	server.router = router
	return server, nil
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
