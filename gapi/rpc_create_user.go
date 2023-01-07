package gapi

import (
	"context"
	db "github.com/hmoodallahma/123bank/db/sqlc"
	"github.com/hmoodallahma/123bank/pb"
	"github.com/hmoodallahma/123bank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context,
	req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"error hashing password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	// todo: refactor --> encapsulate to a util function
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists,
					"username already exists, %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal,
			"failed to create user: %s", err)
	}
	res := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	// todo: set res.User.CreatedAt
	// todo: set res.User.PasswordUpdatedAt
	return res, nil

}
