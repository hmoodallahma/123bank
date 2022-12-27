package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/hmoodallahma/123bank/db/sqlc"
)

//Servers HTTP requests for banking service
type Server struct {
	store *db.Store // access db
	router *gin.Engine // handles routing
}

//creates a new server instance
func NewServer(store *db.Store) *Server{
	server := &Server{store: store}
	router := gin.Default()

	// add routes to router
	router.POST("/account", server.createAccount) // neeed to implement
	router.GET("/account/:id", server.getAccount) // neeed to implement
	router.PUT("/account/:id", server.updateAccount) // neeed to implement
	router.DELETE("/account/:id", server.deleteAccount) // neeed to implement
	router.GET("/accounts", server.listAccounts) // neeed to implement
	server.router = router
	return server
}

// start the server on a specific address
func (server *Server) Start(address string) error {
	// router is private and cannot be accessed outside of pacakge
	return server.router.Run(address)
}

func errorResponse(err error) gin.H{
	return gin.H{
		"error": err.Error(),
	}
}