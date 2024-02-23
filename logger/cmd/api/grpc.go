package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/DpodDani/logger/data"
	logs "github.com/DpodDani/logger/proto"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed to log entry"}
		return res, err
	}

	res := &logs.LogResponse{Result: "successfully logged entry"}
	return res, nil
}

func (app *Config) gRPCListen() {
	listener, err := net.Listen("tpc", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC :%v\n", err)
	}

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})
	log.Printf("gRPC server started and listening on: %s\n", gRPCPort)

	if err = s.Serve(listener); err != nil {
		log.Fatalf("Failed to listen for gRPC :%v\n", err)
	}
}
