package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/hmoodallahma/123bank/db/sqlc"
)

// Servers HTTP requests for banking service
type Server struct {
	store  db.Store    // access db
	router *gin.Engine // handles routing
}

// creates a new server instance
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
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
	server.router = router
	return server
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
