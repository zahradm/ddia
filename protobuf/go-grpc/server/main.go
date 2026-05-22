package server
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"ddia/protobuf/go-grpc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var usersDb = make(map[int32]*user.User)

type server struct {
	user.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	if foundUser, exists := usersDb[req.Id]; exists {
		return &user.UserResponse{User: foundUser}, nil
	}
	return nil, status.Errorf(codes.NotFound, "user not found")
}

func (s *server) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	newUser := req.GetUser()
	if _, exists := usersDb[newUser.Id]; exists {
		return nil, status.Errorf(codes.AlreadyExists, "user with this ID already exists")
	}
	usersDb[newUser.Id] = newUser
	fmt.Printf("Created user: %v\n", newUser)
	return &user.CreateUserResponse{User: newUser}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
