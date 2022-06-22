package main

import (
	"context"
	"log"
	"time"
	"user-service/data"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct{}

// RPCPayload is the type for data we receive from RPC
type RPCPayload struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
	Active    string
}

// WriteUser writes our payload to mongo
func (r *RPCServer) WriteUser(payload RPCPayload, resp *string) error {
	collection := client.Database("users").Collection("users")
	_, err := collection.InsertOne(context.TODO(), data.User{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Password:  payload.Password,
		Active:    payload.Active,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	// resp is the message sent back to the RPC caller
	*resp = "Processed payload via RPC:" + payload.FirstName + " " + payload.LastName
	log.Println("Processed payload via RPC", err)

	return nil
}
