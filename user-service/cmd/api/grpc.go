package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"user-service/data"
	"user-service/users"
)

type UserServer struct {
	users.UnimplementedUserServiceServer
	Models data.Models
}

func (l *UserServer) WriteUser(ctx context.Context, req *users.UserRequest) (*users.UserResponse, error) {
	input := req.GetUser()

	// write the user
	user := data.User{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Password:  input.Password,
		Active:    input.Active,
	}

	err := l.Models.User.Insert(user)
	if err != nil {
		res := &users.UserResponse{Result: "failed"}
		return res, err
	}

	// return response
	res := &users.UserResponse{Result: "inserted"}
	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Faililed to Listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	users.RegisterUserServiceServer(s, &UserServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Listen for gRPC: %v", err)
	}

}
