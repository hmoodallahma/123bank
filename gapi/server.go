package gapi

import (
	"fmt"
	db "github.com/hmoodallahma/123bank/db/sqlc"
	"github.com/hmoodallahma/123bank/pb"
	"github.com/hmoodallahma/123bank/token"
	"github.com/hmoodallahma/123bank/util"
)

// Servers grpc requests for banking service
type Server struct {
	pb.UnimplementedBankServer
	config     util.Config
	store      db.Store    // access db
	tokenMaker token.Maker // make auth tokens
}

// creates a new grpc server instance
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
	return server, nil
}
