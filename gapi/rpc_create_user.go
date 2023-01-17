package gapi

import (
	"context"
	db "github.com/hmoodallahma/123bank/db/sqlc"
	"github.com/hmoodallahma/123bank/pb"
	"github.com/hmoodallahma/123bank/util"
	"github.com/hmoodallahma/123bank/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context,
	req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
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

func validateCreateUserRequest(req *pb.CreateUserRequest) (
	violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
