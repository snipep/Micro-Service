package main

import (
	"context"
	"log"
	"time"

	"github.com/snipep/logger-service/data"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error writing to mongo:", err)
	}

	*resp = "Processed Payload via RPC:" + payload.Name


	return nil
}