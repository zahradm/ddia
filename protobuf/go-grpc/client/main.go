package main

import (
	"context"
	"log"
	"time"

	"ddia/protobuf/go-grpc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
