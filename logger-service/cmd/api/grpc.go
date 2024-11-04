package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/snipep/logger-service/data"
	"github.com/snipep/logger-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (* logs.LogResponse, error) {
	input := req.GetLogEntry()

	//write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed to write log",
		}
		return res, err
	}

	//return response
	res := &logs.LogResponse{
		Result: "logged via grpc",
	}
	return res, err
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s,  &LogServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", gRpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}