# Protobuf with gRPC in Go

This example demonstrates how to use Protocol Buffers with gRPC in Go.

## 1. Define the Protobuf Schema

We will define a `UserService` in our `.proto` file. Create a file named `user.proto` with the following content:

```proto
syntax = "proto3";

package user;

option go_package = ".;user";

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
}

message UserRequest {
  int32 id = 1;
}

message UserResponse {
  User user = 1;
}

message CreateUserRequest {
  User user = 1;
}

message CreateUserResponse {
  User user = 1;
}

service UserService {
  rpc GetUser(UserRequest) returns (UserResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
}
```

## 2. Generate Go Code

First, install the `protoc` compiler and the Go plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Then, generate the Go code:

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto
```

This will create `user.pb.go` and `user_grpc.pb.go`.

## 3. Create the gRPC Server

Create a file named `server/main.go`:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ddia/protobuf/go-grpc/user"
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
```

## 4. Create the gRPC Client

Create a file named `client/main.go`:

```go
package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"ddia/protobuf/go-grpc/user"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := user.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create User
	newUser := &user.User{Id: 1, Name: "gRPC Go User", Email: "go-grpc@example.com"}
	createRes, err := c.CreateUser(ctx, &user.CreateUserRequest{User: newUser})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
	log.Printf("Created User: %v", createRes.GetUser())

	// Get User
	res, err := c.GetUser(ctx, &user.UserRequest{Id: 1})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("Got User: %s", res.GetUser())
}
```

## 5. Setup Go Module

```bash
go mod init ddia/protobuf/go-grpc
go mod tidy
```

## 6. Run the Example

Start the server:
```bash
go run server/main.go
```

In another terminal, run the client:
```bash
go run client/main.go
```
