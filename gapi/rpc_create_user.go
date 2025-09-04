package gapi

import (
	"context"

	db "github.com/hykura1501/simple_bank/db/sqlc"
	"github.com/hykura1501/simple_bank/pb"
	"github.com/hykura1501/simple_bank/util"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				return nil, status.Errorf(codes.AlreadyExists, "username already existes: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "fail to create user: %s", err)
	}
	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
