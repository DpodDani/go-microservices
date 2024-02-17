package main

import (
	"context"
	"log"
	"time"

	"github.com/DpodDani/logger/data"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("Failed to insert log entry into MongoDB: %s\n", err)
		return err
	}

	*resp = "Processed payload via RPC:" + payload.Data
	return nil
}
