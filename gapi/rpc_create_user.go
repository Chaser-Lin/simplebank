package gapi

import (
	db "SimpleBank/db/sqlc"
	"SimpleBank/pb"
	"SimpleBank/util"
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to hash password: %s", err)
	}

	//fmt.Println(hashedPassword)
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	err = server.store.CreateUser(ctx, arg)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok { //转换为mysql类型的error进一步判断错误类型
			switch mysqlErr.Number {
			case 1062: //1062:Duplicate (username)
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	user, err := server.store.GetUser(ctx, arg.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to get user: %s", err)
	}

	fmt.Println(user)
	rsp := &pb.CreateUserResponse{
		//User: convertUser(user),
	}
	fmt.Println(rsp)

	return rsp, nil
}
